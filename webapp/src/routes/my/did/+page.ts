import { getKeyringFromLocalStorage } from '@/keypairoom/keypair';
import { loadFeatureFlags } from '@/features';
import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';
import { redirect } from '@/i18n';

export const load = async ({ fetch }) => {
	const { DID, KEYPAIROOM } = await loadFeatureFlags(fetch);
	if (!KEYPAIROOM && !DID) error(404);

	const keyring = getKeyringFromLocalStorage();
	if (!keyring) redirect('/keypairoom/regenerate');

	const { did } = await pb.send<{ did: JSON }>('/api/did', {});
	return { did };
};
