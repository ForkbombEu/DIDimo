<script lang="ts">
	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import type { Snippet } from 'svelte';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '@/components/ui/utils.js';
	import A from '@/components/ui-custom/a.svelte';

	let {
		ref = $bindable(null),
		class: className,
		href = undefined,
		child,
		children,
		...restProps
	}: WithElementRef<HTMLAnchorAttributes> & {
		child?: Snippet<[{ props: HTMLAnchorAttributes }]>;
	} = $props();

	const attrs = $derived({
		class: cn('hover:text-foreground transition-colors', className),
		href,
		...restProps
	});
</script>

{#if child}
	{@render child({ props: attrs })}
{:else}
	<A bind:this={ref} {...attrs}>
		{@render children?.()}
	</A>
{/if}
