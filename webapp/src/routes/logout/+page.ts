import { browser } from '$app/environment';
import { currentUser, pb } from '@/pocketbase';
import { redirect } from '@/i18n';

export const load = async ({ url }) => {
	if (!browser) return;
	localStorage.clear();
	pb.authStore.clear();
	currentUser.set(null);
	redirect('/', url);
};
