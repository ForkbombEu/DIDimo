<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { page } from '$app/state';
	import Icon from '@/components/ui-custom/icon.svelte';
	import { cn } from '@/components/ui/utils';
	import type { LinkWithIcon } from '../types';
	import { localizeHref } from '@/i18n';
	let { href, icon, title, ...rest }: LinkWithIcon = $props();

	//

	const isActive = $derived(page.url.pathname == href);

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

<a href={href ? localizeHref(href) : undefined} {...rest} role="tab" class={classes}>
	{#if icon}
		<Icon src={icon} mr></Icon>
	{/if}
	{title}
</a>
