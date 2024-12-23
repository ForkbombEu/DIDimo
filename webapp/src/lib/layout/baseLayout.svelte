<script lang="ts">
	import Footer from '$lib/layout/footer.svelte';
	import Topbar from '$lib/layout/topbar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="flex min-h-screen flex-col">
	<Topbar />
	<div class="grow">
		{@render children()}
	</div>
	<Footer />
</div>
