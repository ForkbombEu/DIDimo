// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { pb } from '@/pocketbase';
import type { OrganizationInfoResponse } from '@/pocketbase/types';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const userId = pb.authStore.record?.id;
	if (!userId) error(500);

	let organizationInfo: OrganizationInfoResponse | undefined = undefined;
	try {
		organizationInfo = await pb
			.collection('organization_info')
			.getFirstListItem(`owner = '${userId}'`, { fetch });
	} catch (e) {
		console.log(e);
	}

	return {
		organizationInfo
	};
};
