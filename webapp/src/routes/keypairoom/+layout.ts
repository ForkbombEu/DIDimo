import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	if (!(await loadFeatureFlags(fetch)).KEYPAIROOM) error(404);
};
