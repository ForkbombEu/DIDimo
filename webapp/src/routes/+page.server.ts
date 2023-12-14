import { superValidate } from 'sveltekit-superforms/client';
import { schema } from './_lib/index.js';
import { fail } from '@sveltejs/kit';
import { join } from 'node:path';
import credentialIssuerSchema from '$lib/openid-vc-typescript-json-schema/openid-credential-issuer/schema.json';
import { pb } from '$lib/pocketbase/index.js';

import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import {
	Collections,
	type CredentialIssuersRecord,
	type CredentialIssuersResponse,
	CredentialIssuersScansErrorOptions as Errors
} from '$lib/pocketbase/types.js';

import prependHttp from 'prepend-http';

import { env } from '$env/dynamic/private';

//

export const load = async () => {
	return {
		form: await superValidate(schema)
	};
};

//

export const actions = {
	default: async ({ request, fetch }) => {
		/* Data validation */

		const data = await request.formData();
		const form = await superValidate(data, schema);
		if (!form.valid) return fail(500, { error: 'INVALID_URL' });

		const url = prependHttp(form.data.url);

		/* Pocketbase record creation */

		let credentialIssuerRecordId: string | undefined = undefined;

		await pb.admins.authWithPassword(env.PB_ADMIN_USER, env.PB_ADMIN_PASS);

		const credentialIssuerRecord = await getCredentialIssuerRecord(url);
		if (credentialIssuerRecord) credentialIssuerRecordId = credentialIssuerRecord.id;
		else {
			const newCredentialIssuerRecord = await createCredentialIssuerRecord(url);
			credentialIssuerRecordId = newCredentialIssuerRecord.id;
		}

		/* JSON analysis */

		try {
			validateJSON(await parseResponseJSON(checkResponseStatus(await getData(url, fetch))));
			return {
				form,
				success: true
			};
		} catch (e) {
			if (e instanceof Error) {
				return fail(404, {
					form,
					error: e.message
				});
			}
		}
	}
};

//

async function getCredentialIssuerRecord(
	url: string
): Promise<CredentialIssuersResponse | undefined> {
	const result = await pb
		.collection(Collections.CredentialIssuers)
		.getFullList<CredentialIssuersResponse>({
			filter: `url = "${url}"`
		});

	if (result.length === 0) return undefined;
	else return result[0];
}

async function createCredentialIssuerRecord(url: string): Promise<CredentialIssuersResponse> {
	return await pb
		.collection(Collections.CredentialIssuers)
		.create({ url } satisfies CredentialIssuersRecord);
}

//

function getCredentialIssuerMetadataFilePath(baseUrl: string): string {
	const PATH = '.well-known/openid-credential-issuer';
	return join(baseUrl, PATH);
}

//

async function getData(url: string, fetchFn = fetch): Promise<Response> {
	try {
		return await fetchFn(getCredentialIssuerMetadataFilePath(url));
	} catch (e) {
		console.log(e);
		throw new Error(Errors.CONNECTION_ERROR);
	}
}

function checkResponseStatus(response: Response): Response {
	console.log(response.statusText);
	if (response.status === 200) return response;
	else {
		if (response.status === 404) throw new Error(Errors.FILE_NOT_FOUND);
		else throw new Error(Errors.CONNECTION_ERROR);
	}
}

async function parseResponseJSON(response: Response): Promise<any> {
	try {
		return await response.json();
	} catch (e) {
		console.log(e);
		throw new Error(Errors.BAD_JSON);
	}
}

function validateJSON(data: any) {
	const ajv = new Ajv({ allErrors: true });
	addFormats(ajv);

	// @ts-ignore
	delete credentialIssuerSchema['$schema'];

	const validate = ajv.compile(credentialIssuerSchema);
	validate(data);

	if (validate.errors) throw new Error(Errors.VALIDATION_ERROR);
}
