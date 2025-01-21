import { pb } from '@/pocketbase/index.js';

export const load = async ({ params, fetch }) => {
	const service = await pb.collection('services').getOne(params.service_id, { fetch });
	return {
		service
	};
};
