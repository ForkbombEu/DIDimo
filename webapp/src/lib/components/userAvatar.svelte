<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { currentUser, pb } from '$lib/pocketbase';
	import { Avatar } from 'flowbite-svelte';
	import BoringAvatar from 'svelte-boring-avatars';

	export let size: 'xs' | 'sm' | 'lg' | 'xl' | 'md' = 'md';

	const sizeToNumber = {
		xs: 24, // to be defined
		sm: 32, // to be defined
		md: 40,
		lg: 80,
		xl: 98 // to be defined
	};

	//@ts-ignore
	$: src = pb.files.getUrl($currentUser, $currentUser?.avatar);
</script>

{#if $currentUser?.avatar}
	<Avatar {size} {src} />
{:else}
	<BoringAvatar variant="beam" size={sizeToNumber[size || 'md']} name={$currentUser?.name} />
{/if}
