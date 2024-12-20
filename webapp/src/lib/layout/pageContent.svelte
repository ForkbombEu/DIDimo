<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
		class?: string;
	};

	const { children, class: className = '' }: Props = $props();
</script>

<div class="bg-muted">
	<div class="mx-auto max-w-screen-xl px-8 py-12 pb-32 {className}">
		{@render children()}
	</div>
</div>
