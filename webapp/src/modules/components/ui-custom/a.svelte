<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import { localizeHref } from '@/i18n/paraglide/runtime';
	import { cn } from '../ui/utils';
	type Props = HTMLAnchorAttributes & {
		children?: Snippet;
	};

	let { children, ...rest }: Props = $props();
</script>

<a
	class={cn('w-full', rest.class)}
	href={rest.href ? localizeHref(rest.href) : undefined}
	{...rest}
>
	{@render children?.()}
</a>
