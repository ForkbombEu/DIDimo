import { pb } from '$lib/pocketbase';
import { Collections, type AuthorizationsExamplesRecord } from '$lib/pocketbase/types';

export const load = async () => {
	try {
		// reference: saas/didimo/admin/pb_migrations/seed_test_data.js
		const data = await pb
			.collection(Collections.AuthorizationsExamples)
			.getFirstListItem<AuthorizationsExamplesRecord>('name="authorizationExample"');
		return {
			data
		};
	} catch (error) {
		return {
			error
		};
	}
};
