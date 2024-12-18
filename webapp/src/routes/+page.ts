import { redirect } from '@/i18n';
import { verifyUser } from '@/auth/verifyUser';

export const load = async ({ url }) => {
	if (await verifyUser()) redirect('/my', url);
};
