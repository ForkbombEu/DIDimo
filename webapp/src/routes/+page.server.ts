import { superValidate } from 'sveltekit-superforms/client';
import { schema } from './_lib/index.js';
import { fail } from '@sveltejs/kit';
import { join } from 'node:path';
import credentialIssuerSchema from '$lib/openid-vc-typescript-json-schema/openid-credential-issuer/schema.json';

import Ajv from 'ajv';
import addFormats from 'ajv-formats';

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

		try {
			const PATH = '.well-known/openid-credential-issuer';
			const req = await fetch(join(form.data.url, PATH));

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
	delete credentialIssuerSchema['$schema'];
	const validate = ajv.compile(credentialIssuerSchema);
	validate(data);
	return validate.errors;
}
