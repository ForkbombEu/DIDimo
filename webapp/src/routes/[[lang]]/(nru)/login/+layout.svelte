<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	import { writable } from 'svelte/store';
	export const currentEmail = writable('');
</script>

<script lang="ts">
	import { page } from '$app/stores';
	import type { Link } from '$lib/utils/types';
	import { A, Button, ButtonGroup, Heading, Hr, P } from 'flowbite-svelte';
	import { featureFlags } from '$lib/features';
	import { m } from '$lib/i18n';

	const modes: Link[] = [
		{
			href: '/login',
			text: m.Email_and_password()
		},
		{
			href: '/login/webauthn',
			text: m.Webauthn()
		}
	];
</script>

<Heading tag="h4">Log in</Heading>
{#if $featureFlags.WEBAUTHN}
	<div class="space-y-2">
		<P size="sm" color="text-gray-500">{m.Choose_your_authentication_method()}</P>
		<ButtonGroup class="w-full">
			{#each modes as { href, text }}
				{@const isActive = $page.url.pathname === href}
				<Button color={isActive ? 'dark' : 'alternative'} {href} class="grow">{text}</Button>
			{/each}
		</ButtonGroup>
	</div>
{/if}

<div class="pt-4">
	<slot />
</div>

<div class="flex flex-col gap-4">
	<Hr hrClass="!m-2" />

	<P color="text-gray-500 dark:text-gray-400" class="self-center" size="sm">
		{m.Dont_have_an_account()}
		<A href="/register">{m.Register_here()}</A>
	</P>
</div>
