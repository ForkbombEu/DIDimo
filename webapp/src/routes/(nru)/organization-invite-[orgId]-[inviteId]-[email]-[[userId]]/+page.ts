import { loadFeatureFlags } from '@/features/index.js';
import { redirect } from '@/i18n/index.js';
import { OrganizationInviteSession } from '@/organizations/invites';
import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';

export const load = async ({ params, url, fetch }) => {
	const featureFlags = await loadFeatureFlags(fetch);
	if (!featureFlags.ORGANIZATIONS) error(404);

	OrganizationInviteSession.start({
		organizationId: params.orgId,
		inviteId: params.inviteId,
		email: decodeURIComponent(params.email),
		userId: params.userId
	});

	if (pb.authStore.token) redirect('/my/organizations', url);
	else {
		if (params.userId) redirect('/login', url);
		else redirect('/register', url);
	}
};
