// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import ProtectedOrgLayout from './protectedOrgLayout.svelte';
import ProtectedOrgUI from './protectedOrgUI.svelte';

export * from './roles';
export { ProtectedOrgLayout, ProtectedOrgUI };

//

import { browser } from '$app/environment';
import { pb } from '$lib/pocketbase';
import type { OrgRole } from './roles';

export async function verifyAuthorizations(organizationId: string) {
	if (!browser) return;
	await pb.send('/verify-org-authorization', {
		method: 'POST',
		body: {
			organizationId,
			url: window.location.href
		},
		requestKey: null
	});
}

export async function verifyRole(organizationId: string, roles: OrgRole[]) {
	await pb.send('/verify-user-role', {
		method: 'POST',
		body: {
			organizationId,
			roles
		},
		requestKey: null
	});
}
