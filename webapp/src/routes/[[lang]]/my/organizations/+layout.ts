import { loadFeatureFlags } from '$lib/features';
import { error } from '@sveltejs/kit';

export const load = async () => {
	if (!(await loadFeatureFlags()).ORGANIZATIONS) error(404);
};
