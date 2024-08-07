<script lang="ts">
	import {
		availableLanguageTags,
		type AvailableLanguageTag,
		languageTag
	} from '$paraglide/runtime';
	import { page } from '$app/stores';
	import { i18n, m } from '.';
	import { Dropdown, Button, DropdownHeader, DropdownItem } from 'flowbite-svelte';
	import { ArrowPath } from 'svelte-heros-v2';

	//

	const languagesDisplay: Record<AvailableLanguageTag, { flag: string; name: string }> = {
		en: { flag: '🇬🇧', name: 'English' },
		it: { flag: '🇮🇹', name: 'Italiano' },
		de: { flag: '🇩🇪', name: 'Deutsch' },
		fr: { flag: '🇫🇷', name: 'Français' },
		da: { flag: '🇩🇰', name: 'Dansk' },
		'pt-BR': { flag: '🇧🇷', name: 'Português' }
	};

	const id = 'language-switcher';
</script>

<Button color="alternative" {id}>
	<div class="flex items-center gap-3">
		<p class="h-6 w-6 text-2xl leading-[1]">{languagesDisplay[languageTag()].flag}</p>
		<p>{languagesDisplay[languageTag()].name}</p>
		<ArrowPath size="18" />
	</div>
</Button>

<Dropdown triggeredBy={`#${id}`} class="w-[215px]">
	<DropdownHeader>
		<span class="block truncate text-xs font-medium tracking-wide text-gray-500">
			{m.Select_language()}
		</span>
	</DropdownHeader>

	{#each availableLanguageTags as lang}
		<DropdownItem href={i18n.route($page.url.pathname)} hreflang={lang}>
			{languagesDisplay[lang].flag}
			&nbsp;
			{languagesDisplay[lang].name}
		</DropdownItem>
	{/each}
</Dropdown>
