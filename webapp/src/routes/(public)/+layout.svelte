<script lang="ts">
	import Topbar from '$lib/layout/topbar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<Topbar></Topbar>
{@render children()}
