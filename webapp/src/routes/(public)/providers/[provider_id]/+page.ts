import { createPocketbaseQueryAgent } from '@/pocketbase/query/index.js';

export const load = async ({ params }) => {
	const provider = await createPocketbaseQueryAgent({
		collection: 'services',
		expand: ['credential_issuers']
	}).getOne(params.provider_id);

	return { provider };
};
