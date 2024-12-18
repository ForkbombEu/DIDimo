<script lang="ts">
	import AuthLayout from '@/auth/authLayout.svelte';
	import LanguageSelect from '@/i18n/languageSelect.svelte';
	import type { Snippet } from 'svelte';

	//

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<AuthLayout>
	{#snippet topbarRight()}
		<LanguageSelect />
	{/snippet}

	{@render children?.()}
</AuthLayout>
