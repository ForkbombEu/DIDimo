import { loadFeatureFlags } from '@/features';
import { error, redirect } from '@sveltejs/kit';

import { deLocalizeUrl } from '@/i18n/paraglide/runtime';
export const ssr = false;

export const load = async ({ fetch, url }) => {
	const flags = await loadFeatureFlags(fetch);
	if (flags.MAINTENANCE) error(503);

	if (flags.DEMO && deLocalizeUrl(url).pathname != '/') redirect(303, '/');
};
