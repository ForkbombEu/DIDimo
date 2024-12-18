import { loadFeatureFlags } from '@/features';
import { getUserPublicKeys } from '@/keypairoom/utils';
import { redirect } from '@/i18n';

export const load = async ({ url, fetch }) => {
	const { KEYPAIROOM, AUTH } = await loadFeatureFlags(fetch);

	if (KEYPAIROOM && AUTH) {
		const publicKeys = await getUserPublicKeys(undefined, fetch);
		if (!publicKeys) redirect('/my/keypairoom', url);
	}
};
