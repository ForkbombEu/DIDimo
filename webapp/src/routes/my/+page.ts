// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { PocketbaseQueryAgent } from '@/pocketbase/query/agent.js';
import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const currentUser = pb.authStore.record;
	if (!currentUser) error(404);

	const userOrgAuthorizations = await new PocketbaseQueryAgent(
		{
			collection: 'orgAuthorizations',
			expand: ['organization'],
			filter: `user = "${currentUser.id}"`
		},
		{ fetch }
	).getFullList();

	if (userOrgAuthorizations.length > 2) error(404);

	return {
		organization: userOrgAuthorizations.at(0)?.expand?.organization
	};
};
