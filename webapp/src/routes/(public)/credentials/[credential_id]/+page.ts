// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { PocketbaseQueryAgent } from '@/pocketbase/query/agent.js';

export const load = async ({ params }) => {
	const credential = await new PocketbaseQueryAgent({
		collection: 'credentials',
		expand: ['credential_issuer']
	}).getOne(params.credential_id);

	return {
		credential
	};
};
