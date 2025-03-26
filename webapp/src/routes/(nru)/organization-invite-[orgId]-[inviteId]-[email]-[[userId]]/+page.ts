import { loadFeatureFlags } from '@/features/index.js';
import { OrganizationInviteSession } from '@/organizations/invites';
import { pb } from '@/pocketbase';
import { error, redirect } from '@sveltejs/kit';

export const load = async ({ params, fetch }) => {
	const featureFlags = await loadFeatureFlags(fetch);
	if (!featureFlags.ORGANIZATIONS) error(404);

	OrganizationInviteSession.start({
		organizationId: params.orgId,
		inviteId: params.inviteId,
		email: decodeURIComponent(params.email),
		userId: params.userId
	});

	if (pb.authStore.token) redirect(303, '/my/organizations');
	else {
		if (params.userId) redirect(303, '/login');
		else redirect(303, '/register');
	}
};
