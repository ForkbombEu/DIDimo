import { loadFeatureFlags } from '$lib/features';
import { error } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export const load = async () => {
	const flags = await loadFeatureFlags();
	if (flags.MAINTENANCE) {
		error(503);
	}
};
