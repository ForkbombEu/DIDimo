import { PocketbaseQueryAgent } from '@/pocketbase/query/agent.js';

export const load = async ({ params }) => {
	const credential = await new PocketbaseQueryAgent({
		collection: 'credentials',
		expand: ['credential_issuer']
	}).getOne(params.credential_id);

	return {
		credential
	};
};
