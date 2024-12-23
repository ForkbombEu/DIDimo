<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
		class?: string;
		contentClass?: string;
	};

	const { children, class: className = '', contentClass = '' }: Props = $props();
</script>

<div class={className}>
	<div class="mx-auto max-w-screen-xl px-8 py-12 {contentClass}">
		{@render children()}
	</div>
</div>
