<script lang="ts">
	import Footer from '$lib/layout/footer.svelte';
	import Topbar from '$lib/layout/topbar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<Topbar />
{@render children()}
<Footer />
