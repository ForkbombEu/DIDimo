import { error } from '@sveltejs/kit';
import { loadFeatureFlags } from '@/features';
import { verifyUser } from '@/auth/verifyUser';
import { redirect } from '@/i18n';

export const load = async ({ url, fetch }) => {
	if (!(await loadFeatureFlags(fetch)).AUTH) error(404);
	if (await verifyUser()) redirect('/my', url);
};
