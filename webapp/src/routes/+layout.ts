import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

import { redirect, deLocalizeUrl } from '@/i18n';

export const ssr = false;

export const load = async ({ fetch, url }) => {
	const flags = await loadFeatureFlags(fetch);
	if (flags.MAINTENANCE) error(503);

	if (flags.DEMO && deLocalizeUrl(url).pathname != '/') redirect('/');
};
