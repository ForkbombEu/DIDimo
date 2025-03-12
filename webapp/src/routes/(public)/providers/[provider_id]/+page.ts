import { PocketbaseQueryAgent } from '@/pocketbase/query/index.js';

export const load = async ({ params }) => {
	const provider = await new PocketbaseQueryAgent({
		collection: 'services',
		expand: ['credential_issuers']
	}).getOne(params.provider_id);

	const providerClaims = await new PocketbaseQueryAgent({
		collection: 'provider_claims',
		filter: `provider = "${params.provider_id}"`
	}).getFullList();

	const hasClaim = providerClaims.length > 0;

	return { provider, hasClaim };
};
