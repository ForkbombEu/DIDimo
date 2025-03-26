import { browser } from '$app/environment';
import { currentUser, pb } from '@/pocketbase';
import { redirect } from '@sveltejs/kit';

export const load = async () => {
	if (!browser) return;
	localStorage.clear();
	pb.authStore.clear();
	currentUser.set(null);
	redirect(303, '/');
};
