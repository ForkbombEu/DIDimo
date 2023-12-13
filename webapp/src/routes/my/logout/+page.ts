import { browser } from '$app/environment';
import { currentUser, pb } from '$lib/pocketbase';
import { redirect } from '@sveltejs/kit';

export const load = async () => {
	if (browser) {
		localStorage.clear();
		pb.authStore.clear();
		currentUser.set(null);
		throw redirect(303, '/');
	}
};
