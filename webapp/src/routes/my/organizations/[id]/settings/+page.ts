import { blockMembersWithoutRoles } from '@/organizations';

export const load = async ({ params, fetch }) => {
	const organizationId = params.id;
	await blockMembersWithoutRoles(organizationId, ['owner'], fetch);
};
