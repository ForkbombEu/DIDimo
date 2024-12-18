<script lang="ts">
	import { page } from '$app/stores';
	import Icon from '@/components/ui-custom/icon.svelte';
	import { cn } from '@/components/ui/utils';
	import { Store } from 'runed';
	import type { LinkWithIcon } from '../types';

	let { href, icon, title, ...rest }: LinkWithIcon = $props();

	//

	const pageState = new Store(page);
	const isActive = $derived(pageState.current.url.pathname == href);

	const classes = $derived(
		cn(
			rest.class,
			'inline-block text-sm font-medium text-center p-4 py-3 border-b-2 flex items-center whitespace-nowrap',
			{
				'border-transparent hover:border-primary/20': !isActive,
				'text-primary border-primary border-b-2': isActive
			}
		)
	);
</script>

<a {...rest} role="tab" class={classes} {href}>
	{#if icon}
		<Icon src={icon} mr></Icon>
	{/if}
	{title}
</a>
