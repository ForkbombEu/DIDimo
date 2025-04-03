<script lang="ts">
	import { initI18n } from '@forkbombeu/temporal-ui';
	import type { Snippet } from 'svelte';
	import '@forkbombeu/temporal-ui/dist/temporal.css';

	type Props = {
		children: Snippet;
	};

	const { children }: Props = $props();
</script>

{#await initI18n() then _}
	{@render children()}
{/await}
