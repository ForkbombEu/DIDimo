import { createPocketbaseQueryRunners } from '@/pocketbase/query/index.js';

export const load = async ({ params }) => {
	const provider = await createPocketbaseQueryRunners({
		collection: 'services',
		expand: ['credential_issuers']
	}).getOne(params.provider_id);

	return { provider };
};
