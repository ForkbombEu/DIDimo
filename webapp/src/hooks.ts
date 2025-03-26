import type { Reroute } from '@sveltejs/kit';
import { deLocalizeUrl } from '@/i18n/paraglide/runtime';

export const reroute: Reroute = (request) => {
	return deLocalizeUrl(request.url).pathname;
};
