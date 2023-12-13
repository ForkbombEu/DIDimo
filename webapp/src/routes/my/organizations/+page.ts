import { pb } from '$lib/pocketbase';
import {
	Collections,
	type OrgAuthorizationsResponse,
	type OrgRolesResponse,
	type OrganizationsResponse
} from '$lib/pocketbase/types';

export const load = async () => {
	type Authorizations = Required<
		OrgAuthorizationsResponse<{
			organization: OrganizationsResponse;
			role: OrgRolesResponse;
		}>
	>;

	const authorizations = await pb
		.collection(Collections.OrgAuthorizations)
		.getFullList<Authorizations>({
			filter: `user = "${pb.authStore.model!.id}"`,
			expand: 'organization,role'
		});

	return { authorizations };
};
