<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import { localizeHref } from '@/i18n/paraglide/runtime';
	import { cn } from '@/components/ui/utils.js';

	type Props = HTMLAnchorAttributes & {
		children?: Snippet;
	};

	let { class: className, href, children, ...rest }: Props = $props();
</script>

<a href={href ? localizeHref(href) : undefined} class={cn(className, 'cursor-pointer')} {...rest}>
	{@render children?.()}
</a>
