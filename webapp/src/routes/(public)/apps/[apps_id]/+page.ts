// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { PocketbaseQueryAgent } from '@/pocketbase/query/index.js';

export const load = async ({ params, fetch }) => {
	const wallet = await new PocketbaseQueryAgent(
		{
			collection: 'wallets'
		},
		{ fetch }
	).getOne(params.apps_id);

	return { wallet };
};
