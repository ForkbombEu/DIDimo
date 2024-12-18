import { createI18n } from '@inlang/paraglide-sveltekit';
import * as runtime from './paraglide/runtime.js';
import * as m from './paraglide/messages.js';

export const i18n = createI18n(runtime);
export { m };
export * from './paraglide/runtime.js';

//

import { goto as sveltekitGoto } from '$app/navigation';
import { redirect as sveltekitRedirect, type Page } from '@sveltejs/kit';
import { get } from 'svelte/store';
import { page } from '$app/stores';
import type { AvailableLanguageTag } from './paraglide/runtime.js';
import { Record } from 'effect';

//

export function resolveRoute(route: string, url: URL) {
	const baseRoute = i18n.route(route);
	return i18n.resolveRoute(baseRoute, i18n.getLanguageFromUrl(url));
}

export function goto(route: string) {
	return sveltekitGoto(resolveRoute(route, get(page).url));
}

export function redirect(route: string, fromUrl: URL, statusCode: RedirectStatusCode = 303) {
	return sveltekitRedirect(statusCode, resolveRoute(route, fromUrl));
}

type RedirectStatusCode = Parameters<typeof sveltekitRedirect>['0'];

//

export const languagesDisplay: Record<AvailableLanguageTag, { flag: string; name: string }> = {
	en: { flag: '🇬🇧', name: 'English' },
	it: { flag: '🇮🇹', name: 'Italiano' },
	de: { flag: '🇩🇪', name: 'Deutsch' },
	fr: { flag: '🇫🇷', name: 'Français' },
	da: { flag: '🇩🇰', name: 'Dansk' },
	'pt-br': { flag: '🇧🇷', name: 'Português' }
};

export function getLanguagesData(page: Page): LanguageData[] {
	const currentLang = i18n.getLanguageFromUrl(page.url);

	return Record.keys(languagesDisplay).map((lang) => ({
		tag: lang,
		href: i18n.route(page.url.pathname),
		hreflang: lang,
		flag: languagesDisplay[lang].flag,
		name: languagesDisplay[lang].name,
		isCurrent: lang == currentLang
	}));
}

export type LanguageData = {
	tag: AvailableLanguageTag;
	href: string;
	hreflang: AvailableLanguageTag;
	flag: string;
	name: string;
	isCurrent: boolean;
};
