<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import { localizeHref } from '@/i18n/paraglide/runtime';
	type Props = HTMLAnchorAttributes & {
		children?: Snippet;
	};

	let { children, ...rest }: Props = $props();
</script>

<a href={rest.href ? localizeHref(rest.href) : undefined} {...rest}>
	{@render children?.()}
</a>

<style lang="postcss">
	a {
		@apply text-primary hover:underline;
	}
</style>
