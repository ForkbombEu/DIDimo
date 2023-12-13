import { loadFeatureFlags } from '$lib/features';

export const ssr = false;

export const load = async () => {
	await loadFeatureFlags();
};
