import { redirect } from '@sveltejs/kit';
import { verifyUser } from '$lib/auth/verifyUser';

export const load = async () => {
	if (await verifyUser()) throw redirect(303, '/my');
};
