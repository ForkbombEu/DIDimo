// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { verifyUser } from '@/auth/verifyUser';
import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

import { browser } from '$app/environment';
import { redirect } from '@/i18n';
import { getKeyringFromLocalStorage, matchPublicAndPrivateKeys } from '@/keypairoom/keypair';
import { getUserPublicKeys, RegenerateKeyringSession } from '@/keypairoom/utils';

import { OrganizationInviteSession } from '@/organizations/invites/index.js';

export const load = async ({ fetch }) => {
	if (!browser) return;
	const featureFlags = await loadFeatureFlags(fetch);

	// Auth

	if (!featureFlags.AUTH) error(404);
	if (!(await verifyUser(fetch))) redirect('/login');

	// Keypairoom

	if (featureFlags.KEYPAIROOM) {
		const publicKeys = await getUserPublicKeys();
		if (!publicKeys) redirect('/keypairoom');

		const keyring = getKeyringFromLocalStorage();
		if (!keyring) {
			RegenerateKeyringSession.start();
			redirect('/keypairoom/regenerate');
		}

		try {
			if (publicKeys && keyring) await matchPublicAndPrivateKeys(publicKeys, keyring);
		} catch {
			RegenerateKeyringSession.start();
			redirect('/keypairoom/regenerate');
		}
	}
	if (featureFlags.KEYPAIROOM && RegenerateKeyringSession.isActive()) {
		RegenerateKeyringSession.end();
	}

	// Organizations

	if (featureFlags.ORGANIZATIONS && OrganizationInviteSession.isActive()) {
		OrganizationInviteSession.end();
		redirect('/my/organizations');
	}
};
