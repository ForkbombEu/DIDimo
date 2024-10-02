// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { CredentialIssuersFeaturesTypeOptions as Feature } from '$lib/pocketbase/types.js';
import Ajv from 'ajv';
import addFormats from 'ajv-formats';
import credentialIssuerSchema from '$lib/openid-vc-typescript-json-schema/openid-credential-issuer/schema.json';

//

type ValidationResult<T> = {
	data: T;
	feature: Feature;
};

export function doesFileExist(response: Response): ValidationResult<Response> {
	if (response.status === 200) return { data: response, feature: Feature.FILE_EXISTS };
	throw new Error(Feature.FILE_EXISTS);
}

export async function getResponseJson(response: Response): Promise<ValidationResult<Response>> {
	try {
		return { data: await response.json(), feature: Feature.VALID_JSON };
	} catch (e) {
		throw new Error(Feature.VALID_JSON);
	}
}

export function validateJson(data: unknown): ValidationResult<unknown> {
	const ajv = new Ajv({ allErrors: true });
	addFormats(ajv);

	// @ts-ignore
	delete credentialIssuerSchema['$schema'];

	const validate = ajv.compile(credentialIssuerSchema);
	validate(data);

	if (validate.errors) throw new Error(Feature.SCHEMA_COMPLIANT);
	return { data, feature: Feature.SCHEMA_COMPLIANT };
}
