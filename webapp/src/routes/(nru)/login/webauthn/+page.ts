import { loadFeatureFlags } from '@/features';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const { WEBAUTHN } = await loadFeatureFlags(fetch);
	if (!WEBAUTHN) {
		error(404);
	}
};
