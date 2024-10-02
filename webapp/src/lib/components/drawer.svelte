<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Drawer, Heading } from 'flowbite-svelte';
	import { sineIn } from 'svelte/easing';
	import PortalWrapper from './portalWrapper.svelte';
	import IconButton from './iconButton.svelte';

	export let hidden = true;
	export let placement: 'left' | 'right' = 'left';
	export let width = 'w-96';
	export let backdrop = true;
	export let closeOnClickOutside = true;
	export let hideTopbar = false;
	export let title: string | undefined = undefined;

	$: transitionDirection = placement == 'right' ? 1 : -1;
	$: transitionParams = {
		x: 320 * transitionDirection,
		duration: 200,
		easing: sineIn
	};
</script>

<PortalWrapper>
	<Drawer
		bind:hidden
		transitionType="fly"
		{backdrop}
		{transitionParams}
		activateClickOutside={closeOnClickOutside}
		{placement}
		class="flex flex-col !p-0"
		{width}
	>
		{#if !hideTopbar}
			<div class="flex justify-between items-center p-2 pl-4 border-b">
				<Heading tag="h6">
					{#if title}
						{title}
					{/if}
				</Heading>
				<IconButton
					on:click={() => {
						hidden = true;
					}}
				/>
			</div>
		{/if}
		<slot />
	</Drawer>
</PortalWrapper>
