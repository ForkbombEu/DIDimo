<script lang="ts" module>
	export type LanguageSelectTriggerSnippetProps = {
		icon: IconComponent;
		text: string;
		language: LanguageData;
	};
</script>

<script lang="ts">
	import { page } from '$app/state';
	import { getLanguagesData, m, type LanguageData } from '.';
	import { Languages } from 'lucide-svelte';
	import { getLocale } from './paraglide/runtime';
	import type { Snippet } from 'svelte';
	import type { IconComponent } from '@/components/types';

	type Props = {
		languages: Snippet<[{ languages: LanguageData[] }]>;
		contentClass?: string;
		trigger: Snippet<[LanguageSelectTriggerSnippetProps]>;
	};

	const { trigger, languages: languagesSnippet }: Props = $props();

	const languages = $derived(getLanguagesData(page));
	const currentLanguage = $derived(languages.find((l) => l.tag == getLocale())!);
</script>

{@render trigger({
	icon: Languages,
	text: m.Select_language(),
	language: currentLanguage
})}

{@render languagesSnippet({ languages })}
