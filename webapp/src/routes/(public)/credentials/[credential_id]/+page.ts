import { pb } from '@/pocketbase';

export const load = async ({ params }) => {
	const credential = await pb.collection('credentials').getOne(params.credential_id);

	return {
		credential
	};
};
