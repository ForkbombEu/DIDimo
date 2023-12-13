import { loadFeatureFlags } from '$lib/features';
import { error } from '@sveltejs/kit';

export const ssr = false;

export const load = async () => {
	const flags = await loadFeatureFlags();
	if (flags.MAINTENANCE) {
		throw error(503);
	}
};
