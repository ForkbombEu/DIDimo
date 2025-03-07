import { pb } from '@/pocketbase/index.js';

export const load = async ({ params }) => {
	const issuer = await pb.collection('credential_issuers').getOne(params.issuer_id);
	return { issuer };
};
