import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const record = (
		await pb.collection('z_test_collection').getFullList({ fetch, requestKey: null })
	).at(0);
	if (!record) error(404);
	return { record };
};
