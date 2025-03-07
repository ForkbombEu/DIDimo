import type { CredentialConfiguration } from '$lib/types/openid.js';
import { pb } from '@/pocketbase';
import type { CredentialsResponse } from '@/pocketbase/types/index.generated.js';

export const load = async ({ params }) => {
	const credential = await pb
		.collection('credentials')
		.getOne<CredentialsResponse<CredentialConfiguration>>(params.credential_id);

	return {
		credential
	};
};
