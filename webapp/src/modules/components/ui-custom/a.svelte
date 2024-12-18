<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';

	type Props = HTMLAnchorAttributes & {
		children?: Snippet;
	};

	let { children, ...rest }: Props = $props();
</script>

<a {...rest}>
	{@render children?.()}
</a>

<style lang="postcss">
	a {
		@apply text-primary hover:underline;
	}
</style>
