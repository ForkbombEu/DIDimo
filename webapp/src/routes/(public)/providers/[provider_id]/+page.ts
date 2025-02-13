import { PocketbaseQuery } from '@/pocketbase/query/index.js';

export const load = async ({ params }) => {
	const provider = await new PocketbaseQuery('services', {
		expand: ['credential_issuers']
	}).getOne(params.provider_id);

	return { provider };
};
