import { pb } from '@/pocketbase';
import type { OrgRole } from '.';
import { type OrgAuthorizationsResponse, type OrgRolesResponse } from '@/pocketbase/types';
import { Option as O } from 'effect';

export async function getUserRole(organizationId: string, userId: string): Promise<OrgRole> {
	type AuthorizationWithRole = OrgAuthorizationsResponse<{ role: OrgRolesResponse }>;

	const orgAuthorization = await pb
		.collection('orgAuthorizations')
		.getFirstListItem<AuthorizationWithRole>(
			`organization.id = '${organizationId}' && user.id = '${userId}'`,
			{
				expand: 'role',
				requestKey: null
			}
		);

	return O.fromNullable(orgAuthorization.expand?.role.name).pipe(O.getOrThrow) as OrgRole;
}
