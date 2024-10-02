// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { superValidate } from 'sveltekit-superforms/client';
import { schema } from './_lib/schema.js';

import { fail } from '@sveltejs/kit';
import { join } from 'node:path';
import prependHttp from 'prepend-http';

import { pb } from '$lib/pocketbase/index.js';
import {
	Collections,
	type CredentialIssuersRecord,
	type CredentialIssuersResponse,
	CredentialIssuersFeaturesTypeOptions as Feature,
	type CredentialIssuersFeaturesRecord
} from '$lib/pocketbase/types.js';

import { doesFileExist, getResponseJson, validateJson } from './_lib/credentialIssuerValidators.js';

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

		try {
			await pb.admins.authWithPassword(env.PB_ADMIN_USER, env.PB_ADMIN_PASS);
			const credentialIssuerRecord = await getOrCreateCredentialIssuerRecord(url);
			const reportRecord = await createCredentialIssuerReportRecord(credentialIssuerRecord.id);

			const response = await getData(url, fetch);
			return {
				form,
				features: await runValidators(response, reportRecord.id)
			};
		} catch (e) {
			if (e instanceof ConnectionError) {
				return fail(404, {
					form,
					connectionError: true
				});
			} else if (e instanceof Error) {
				return fail(500, {
					form,
					error: e.message
				});
			} else {
				return fail(500, {
					form,
					error: JSON.stringify(e)
				});
			}
		}
	}
};

/* Pocketbase operations */

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

async function getOrCreateCredentialIssuerRecord(url: string) {
	const credentialIssuerRecord = await getCredentialIssuerRecord(url);
	if (!credentialIssuerRecord) return await createCredentialIssuerRecord(url);
	return credentialIssuerRecord;
}

async function createCredentialIssuerReportRecord(credentialIssuerRecordId: string) {
	return await pb
		.collection(Collections.CredentialIssuersReports)
		.create<CredentialIssuersResponse>({
			credential_issuer: credentialIssuerRecordId
		});
}

async function createCredentialIssuerFeatureRecord(
	credentialIssuerReportId: string,
	feature: Feature
) {
	await pb.collection(Collections.CredentialIssuersFeatures).create({
		report: credentialIssuerReportId,
		type: feature
	} satisfies CredentialIssuersFeaturesRecord);
}

/* Data fetch operations */

function getCredentialIssuerMetadataFilePath(baseUrl: string): string {
	const PATH = '.well-known/openid-credential-issuer';
	return join(baseUrl, PATH);
}

async function getData(url: string, fetchFn = fetch): Promise<Response> {
	try {
		return await fetchFn(getCredentialIssuerMetadataFilePath(url));
	} catch (e) {
		console.log(e);
		throw new ConnectionError();
	}
}

class ConnectionError extends Error {
	constructor() {
		super('connection_error');
		this.name = this.constructor.name;
		if (Error.captureStackTrace) {
			Error.captureStackTrace(this, this.constructor);
		}
	}
}

/* Validation operations */

async function runValidators(response: Response, reportRecordId: string): Promise<Array<Feature>> {
	const features: Feature[] = [];

	try {
		const validResponse = doesFileExist(response);
		await createCredentialIssuerFeatureRecord(reportRecordId, validResponse.feature);
		features.push(validResponse.feature);

		const json = await getResponseJson(validResponse.data);
		await createCredentialIssuerFeatureRecord(reportRecordId, json.feature);
		features.push(json.feature);

		const validJson = validateJson(json.data);
		await createCredentialIssuerFeatureRecord(reportRecordId, validJson.feature);
		features.push(validJson.feature);
	} catch (e) {
		console.log(e);
	}

	return features;
}
