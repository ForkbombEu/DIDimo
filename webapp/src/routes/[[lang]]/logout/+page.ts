import { browser } from '$app/environment';
import { currentUser, pb } from '$lib/pocketbase';
import { redirect } from '$lib/i18n';

export const load = async ({ url }) => {
	if (browser) {
		localStorage.clear();
		pb.authStore.clear();
		currentUser.set(null);
		redirect('/', url);
	}
};
