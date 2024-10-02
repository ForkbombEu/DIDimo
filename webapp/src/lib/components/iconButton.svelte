<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import { XMark } from 'svelte-heros-v2';
	import type { IconComponent } from '$lib/utils/types';

	export let icon: IconComponent = XMark;
	export let border = false;
	export let size: keyof typeof sizes = 'xs';
	export let href: string | undefined = undefined;

	const sizes = {
		xs: {
			iconSize: 20,
			padding: '!p-1'
		},
		sm: {
			iconSize: 20,
			padding: '!p-2'
		}
	};

	$: props = sizes[size];
</script>

<Button
	class={`${props.padding} ${!border ? '!border-transparent' : ''} !shrink-0`}
	color="alternative"
	on:click
	{href}
>
	<svelte:component this={icon} size={props.iconSize.toString()} />
</Button>
