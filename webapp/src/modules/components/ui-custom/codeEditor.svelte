<script lang="ts">
	import CodeMirror from 'svelte-codemirror-editor';
	import { json } from '@codemirror/lang-json';
	import { dracula } from 'thememirror';

	//

	type LanguageSupport = ReturnType<typeof json>;
	type Extension = typeof dracula;

	const langs = {
		json
	};

	const themes = {
		dracula
	};

	//

	type Props = {
		minHeight?: number | null;
		maxHeight?: number | null;
		value?: string | null | undefined;
		lang: keyof typeof langs | LanguageSupport;
		theme?: keyof typeof themes | Extension;
		class?: string;
	};

	let {
		lang,
		minHeight = 100,
		maxHeight,
		theme = 'dracula',
		value = $bindable(),
		class: className = ''
	}: Props = $props();

	//

	const languageSupport: LanguageSupport | null = $derived.by(() => {
		if (typeof lang == 'string') {
			if (lang in langs) return langs[lang]();
			else return null;
		} else {
			return lang;
		}
	});

	const themeExtension: Extension | null = $derived.by(() => {
		if (typeof theme == 'string') {
			if (theme in themes) return themes[theme];
			else return null;
		} else {
			return theme;
		}
	});

	const styles = $derived.by(() => {
		const baseStyles = {
			'&': { minHeight: 'none', maxHeight: 'none' },
			'.cm-scroller': { overflow: 'auto' }
		};
		if (minHeight) baseStyles['&'].minHeight = `${minHeight}px`;
		if (maxHeight) baseStyles['&'].maxHeight = `${maxHeight}px`;
		return baseStyles;
	});
</script>

<CodeMirror
	lang={languageSupport}
	theme={themeExtension}
	class="overflow-hidden rounded-lg {className}"
	{styles}
	bind:value
/>
