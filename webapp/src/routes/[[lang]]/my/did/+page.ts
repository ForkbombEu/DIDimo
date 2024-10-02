// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { getKeyringFromLocalStorage } from '$lib/keypairoom/keypair';
import { loadFeatureFlags } from '$lib/features';
import { pb } from '$lib/pocketbase';
import { error } from '@sveltejs/kit';
import { redirect } from '$lib/i18n';

export const load = async ({ url }) => {
	const { DID, KEYPAIROOM } = await loadFeatureFlags();
	if (!KEYPAIROOM && !DID) error(404);

	const keyring = getKeyringFromLocalStorage();
	if (!keyring) redirect('/keypairoom/regenerate', url);

	const { did } = await pb.send<{ did: JSON }>('/api/did', {});
	return { did };
};
