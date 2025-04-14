// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { PocketbaseQueryAgent } from '@/pocketbase/query/index.js';

export const load = async ({ params, fetch }) => {
	const service = await new PocketbaseQueryAgent(
		{
			collection: 'credential_issuers',
			expand: ['credentials_via_credential_issuer']
		},
		{ fetch }
	).getOne(params.service_id);

	return { service };
};
