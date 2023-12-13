<script lang="ts" context="module">
	import { writable } from 'svelte/store';
	export const currentEmail = writable('');
</script>

<script lang="ts">
	import { page } from '$app/stores';
	import type { Link } from '$lib/utils/types';
	import { A, Button, ButtonGroup, Heading, Hr, P } from 'flowbite-svelte';
	import { featureFlags } from '$lib/features';

	const modes: Link[] = [
		{
			href: '/login',
			text: 'Email and password'
		},
		{
			href: '/login/webauthn',
			text: 'Webauthn'
		}
	];
</script>

<Heading tag="h4">Log in</Heading>
{#if $featureFlags.WEBAUTHN}
	<div class="space-y-2">
		<P size="sm" color="text-gray-500">Choose your authentication method</P>
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

<div class="flex flex-col items-center gap-4">
	<Hr />
	<P color="text-gray-500 dark:text-gray-400" size="sm">
		Don't have an account?
		<A href="/register">Register here</A>
	</P>
</div>
