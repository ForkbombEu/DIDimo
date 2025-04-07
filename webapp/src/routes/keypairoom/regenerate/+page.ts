// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { loadFeatureFlags } from '@/features';
import { getUserPublicKeys } from '@/keypairoom/utils';
import { redirect } from '@/i18n';

export const load = async ({ fetch }) => {
	const { KEYPAIROOM, AUTH } = await loadFeatureFlags(fetch);

	if (KEYPAIROOM && AUTH) {
		const publicKeys = await getUserPublicKeys(undefined, fetch);
		if (!publicKeys) redirect('/my/keypairoom');
	}
};
