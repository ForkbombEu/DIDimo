import { redirect } from '$lib/i18n';
import { verifyUser } from '$lib/auth/verifyUser';

export const load = async ({ url }) => {
	if (await verifyUser()) redirect('/my', url);
};
