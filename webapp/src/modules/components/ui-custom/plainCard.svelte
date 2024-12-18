<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		class?: string;
		children?: Snippet;
	}

	let { class: className = '', children }: Props = $props();
</script>

<div class="rounded-md border px-4 py-3 text-sm {className}">
	{@render children?.()}
</div>
