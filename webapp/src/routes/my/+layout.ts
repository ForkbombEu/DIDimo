import { verifyUser } from '@/auth/verifyUser';
import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { getKeyringFromLocalStorage, matchPublicAndPrivateKeys } from '@/keypairoom/keypair';
import { getUserPublicKeys, RegenerateKeyringSession } from '@/keypairoom/utils';

import { OrganizationInviteSession } from '@/organizations/invites/index.js';

export const load = async ({ fetch }) => {
	if (!browser) return;
	const featureFlags = await loadFeatureFlags(fetch);

	// Auth

	if (!featureFlags.AUTH) error(404);
	if (!(await verifyUser(fetch))) redirect(303, '/login');

	// Keypairoom

	if (featureFlags.KEYPAIROOM) {
		const publicKeys = await getUserPublicKeys();
		if (!publicKeys) redirect(303, '/keypairoom');

		const keyring = getKeyringFromLocalStorage();
		if (!keyring) {
			RegenerateKeyringSession.start();
			redirect(303, '/keypairoom/regenerate');
		}

		try {
			if (publicKeys && keyring) await matchPublicAndPrivateKeys(publicKeys, keyring);
		} catch {
			RegenerateKeyringSession.start();
			redirect(303, '/keypairoom/regenerate');
		}
	}
	if (featureFlags.KEYPAIROOM && RegenerateKeyringSession.isActive()) {
		RegenerateKeyringSession.end();
	}

	// Organizations

	if (featureFlags.ORGANIZATIONS && OrganizationInviteSession.isActive()) {
		OrganizationInviteSession.end();
		redirect(303, '/my/organizations');
	}
};
