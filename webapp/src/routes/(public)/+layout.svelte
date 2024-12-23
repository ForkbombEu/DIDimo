<script lang="ts">
	import BaseLayout from '$lib/layout/baseLayout.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<BaseLayout>
	{@render children()}
</BaseLayout>
