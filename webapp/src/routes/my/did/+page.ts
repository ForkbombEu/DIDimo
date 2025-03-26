import { getKeyringFromLocalStorage } from '@/keypairoom/keypair';
import { loadFeatureFlags } from '@/features';
import { pb } from '@/pocketbase';
import { error, redirect } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const { DID, KEYPAIROOM } = await loadFeatureFlags(fetch);
	if (!KEYPAIROOM && !DID) error(404);

	const keyring = getKeyringFromLocalStorage();
	if (!keyring) redirect(303, '/keypairoom/regenerate');

	const { did } = await pb.send<{ did: JSON }>('/api/did', {});
	return { did };
};
