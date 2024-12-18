<script lang="ts" module>
	export type LanguageSelectTriggerSnippetProps = {
		icon: IconComponent;
		text: string;
		language: LanguageData;
	};
</script>

<script lang="ts">
	import { page } from '$app/stores';
	import { getLanguagesData, m, type LanguageData } from '.';
	import { Languages } from 'lucide-svelte';
	import { languageTag } from '.';
	import type { Snippet } from 'svelte';
	import type { IconComponent } from '@/components/types';
	import { Store } from 'runed';

	type Props = {
		languages: Snippet<[{ languages: LanguageData[] }]>;
		contentClass?: string;
		trigger: Snippet<[LanguageSelectTriggerSnippetProps]>;
	};

	const { trigger, languages: languagesSnippet }: Props = $props();

	const pageState = new Store(page);
	const languages = $derived(getLanguagesData(pageState.current));
	const currentLanguage = $derived(languages.find((l) => l.tag == languageTag())!);
</script>

{@render trigger({
	icon: Languages,
	text: m.Select_language(),
	language: currentLanguage
})}

{@render languagesSnippet({ languages })}
