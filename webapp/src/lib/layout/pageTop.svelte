<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
		contentClass?: string;
	};

	const { children, contentClass = '' }: Props = $props();
</script>

<div class="bg-card border-b-primary border-b-2">
	<div class="mx-auto max-w-screen-xl space-y-8 px-8 py-12 pt-20 {contentClass}">
		{@render children()}
	</div>
</div>
