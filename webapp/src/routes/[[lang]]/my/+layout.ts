// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { browser } from '$app/environment';
import { verifyUser } from '$lib/auth/verifyUser';
import { loadFeatureFlags } from '$lib/features';
import { getUserPublicKeys } from '$lib/keypairoom/utils.js';
import { getKeyringFromLocalStorage, matchPublicAndPrivateKeys } from '$lib/keypairoom/keypair';
import { missingKeyringParam, welcomeSearchParamKey } from '$lib/utils/constants.js';
import { error } from '@sveltejs/kit';
import { redirect } from '$lib/i18n';

//

export const load = async ({ url }) => {
	const features = await loadFeatureFlags();

	if (!features.AUTH) error(404);
	else {
		const isUserLogged = await verifyUser();
		if (!isUserLogged) redirect('/login', url);
	}

	if (features.KEYPAIROOM) {
		const publicKeys = await getUserPublicKeys();
		if (!publicKeys) {
			redirect(`/keypairoom?${welcomeSearchParamKey}`, url);
			return;
		}

		if (browser) {
			const keyring = getKeyringFromLocalStorage();

			if (!keyring) {
				redirect(`/keypairoom/regenerate?${missingKeyringParam}`, url);
				return;
			}

			try {
				await matchPublicAndPrivateKeys(publicKeys, keyring);
			} catch (e) {
				redirect(`/keypairoom/regenerate?${missingKeyringParam}`, url);
			}
		}
	}
};
