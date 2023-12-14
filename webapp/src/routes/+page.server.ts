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
	type CredentialIssuersResponse
} from '$lib/pocketbase/types.js';

import { env } from '$env/dynamic/private';

export const load = async () => {
	return {
		form: await superValidate(schema)
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const form = await superValidate(data, schema);

		if (!form.valid) return fail(400, { form, message: 'Invalid URL provided' });

		const { url } = form.data;

		try {
			const PATH = '.well-known/openid-credential-issuer';
			const req = await fetch(join(url, PATH));

			await pb.admins.authWithPassword(env.PB_ADMIN_USER, env.PB_ADMIN_PASS);

			let id: string;

			const res = await pb
				.collection(Collections.CredentialIssuers)
				.getFullList<CredentialIssuersResponse>({
					filter: `url = "${url}"`
				});

			if (res.length === 0) {
				const newRecord = await pb
					.collection(Collections.CredentialIssuers)
					.create({ url } satisfies CredentialIssuersRecord);
				id = newRecord.id;
			} else {
				id = res[0].id;
			}

			if (req.status === 404) {
				return fail(404, {
					form,
					message: `Credential Issuer metadata file not found ( ${PATH} )`
				});
			}

			if (req.status !== 200) {
				return fail(404, {
					form,
					message: req.statusText
				});
			}

			const JSON = await req.json();
			const errors = validateJSON(JSON);

			if (errors) {
				return fail(500, {
					form,
					validationErrors: errors
				});
			} else {
				return {
					success: true,
					form
				};
			}
		} catch (e) {
			return fail(404, {
				form,
				message: e.message
			});
		}
	}
};

function validateJSON(data: any) {
	const ajv = new Ajv({ allErrors: true });
	addFormats(ajv);
	delete credentialIssuerSchema['$schema']; // Check type safety
	const validate = ajv.compile(credentialIssuerSchema);
	validate(data);
	return validate.errors;
}
