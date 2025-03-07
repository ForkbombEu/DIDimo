import { PocketbaseQueryAgent } from '@/pocketbase/query/index.js';

export const load = async ({ params }) => {
	const provider = await new PocketbaseQueryAgent({
		collection: 'services',
		expand: ['credential_issuers']
	}).getOne(params.provider_id);

	return { provider };
};
