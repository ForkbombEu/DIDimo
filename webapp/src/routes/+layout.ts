import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

export const ssr = false;

export const load = async ({ fetch }) => {
	const flags = await loadFeatureFlags(fetch);
	if (flags.MAINTENANCE) {
		error(503);
	}
};
