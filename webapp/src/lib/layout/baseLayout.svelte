<script lang="ts">
	import Footer from '$lib/layout/footer.svelte';
	import Topbar from '$lib/layout/topbar.svelte';
	import { featureFlags } from '@/features';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="flex min-h-screen flex-col">
	{#if !$featureFlags.DEMO}
		<Topbar />
	{/if}
	<div class="flex grow flex-col">
		{@render children()}
	</div>
	<Footer />
</div>
