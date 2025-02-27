import { m } from '@/i18n';
import { Cog, Home, Users, Network } from 'lucide-svelte';
import type { OrgRole } from '.';
import type { LinkWithIcon } from '@/components/types';

//

export function createOrganizationLinks(
	organizationId: string,
	userRole: OrgRole = 'member'
): LinkWithIcon[] {
	const base = (path = '') => `/my/organizations/${organizationId}${path}`;

	const links: LinkWithIcon[] = [
		{
			title: m.Home(),
			href: base(),
			icon: Home
		},
		{
			title: m.Services(),
			href: base('/providers'),
			icon: Network
		},
		{
			title: m.Members(),
			href: base('/members'),
			icon: Users
		}
	];

	if (userRole == 'owner') {
		links.push({
			title: m.Settings(),
			href: base('/settings'),
			icon: Cog
		});
	}

	return links;
}
