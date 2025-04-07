// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import type { OpenidCredentialIssuerSchema } from './openid-credential-issuer';

export type { OpenidCredentialIssuerSchema };
export type CredentialConfiguration =
	OpenidCredentialIssuerSchema['credential_configurations_supported'][string];
export type CredentialDefinition = NonNullable<CredentialConfiguration['credential_definition']>;
export type CredentialSubject = NonNullable<CredentialDefinition['credentialSubject']>;
