<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		class?: string;
		children?: Snippet;
	};

	let { class: className = '', children }: Props = $props();
</script>

<div class="grid grid-cols-1 gap-6 md:grid-cols-2 md:gap-8 lg:grid-cols-3 {className}">
	{@render children?.()}
</div>
