import { browser } from '$app/environment';
import { pb } from '$lib/pocketbase';
import { Collections } from '$lib/pocketbase/types';

// Reference: https://github.com/pocketbase/js-sdk/issues/85

export async function verifyUser(): Promise<boolean> {
	if (!browser) throw new Error('verifyUser() must be called from the browser');
	if (!document || !document.cookie) return false;

	pb.authStore.loadFromCookie(document.cookie);
	try {
		await pb.collection(Collections.Users).authRefresh();
	} catch (_) {
		pb.authStore.clear();
		return false;
	}
	return true;
}
