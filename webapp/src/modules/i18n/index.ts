import { Record } from 'effect';
import type { Handle, Page } from '@sveltejs/kit';

import { paraglideMiddleware } from './paraglide/server';
import { locales, localizeHref, getLocale } from '@/i18n/paraglide/runtime';
import * as m from './paraglide/messages.js';

export { m };

export const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request: localizedRequest, locale }) => {
		event.request = localizedRequest;
		return resolve(event, {
			transformPageChunk: ({ html }) => {
				return html.replace('%lang%', locale);
			}
		});
	});

export const languagesDisplay: Record<(typeof locales)[number], { flag: string; name: string }> = {
	en: { flag: '🇬🇧', name: 'English' },
	it: { flag: '🇮🇹', name: 'Italiano' },
	de: { flag: '🇩🇪', name: 'Deutsch' },
	fr: { flag: '🇫🇷', name: 'Français' },
	da: { flag: '🇩🇰', name: 'Dansk' },
	'pt-br': { flag: '🇧🇷', name: 'Português' }
};

export function getLanguagesData(page: Page): LanguageData[] {
	const currentLang = getLocale();

	return Record.keys(languagesDisplay).map((lang) => ({
		tag: lang,
		href: localizeHref(page.url.pathname, { locale: lang }),
		hreflang: lang,
		flag: languagesDisplay[lang].flag,
		name: languagesDisplay[lang].name,
		isCurrent: lang == currentLang
	}));
}

export type LanguageData = {
	tag: (typeof locales)[number];
	href: string;
	hreflang: (typeof locales)[number];
	flag: string;
	name: string;
	isCurrent: boolean;
};
