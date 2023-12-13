import { loadFeatureFlags } from '$lib/features';
import { error } from '@sveltejs/kit';

export const load = async () => {
	const { WEBAUTHN } = await loadFeatureFlags();
	if (!WEBAUTHN) {
		throw error(404);
	}
};
