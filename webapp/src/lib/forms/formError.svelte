<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Alert } from 'flowbite-svelte';
	import { getFormContext } from './form.svelte';
	import { m } from '$lib/i18n';

	const { superform } = getFormContext();
	const { message, allErrors } = superform;

	$: error = $allErrors.at(0);
</script>

{#if error}
	<Alert color="red" dismissable>
		<p class="font-bold">{m.Error()}</p>
		{#if error.messages.length > 0}
			<ul class="mt-1 space-y-2">
				{#each error.messages as message}
					<li class="leading-4">
						{message}
					</li>
				{/each}
			</ul>
		{/if}
	</Alert>
{/if}

{#if $message}
	<Alert color="red" dismissable>
		<p>{$message.message}</p>
		{#if $message.data && Object.keys($message.data).length > 0}
			<ul class="mt-2 space-y-2">
				{#each Object.entries($message.data) as [key, value]}
					<li class="leading-4">
						<span class="font-bold">{key}</span><br />{value.message}
					</li>
				{/each}
			</ul>
		{/if}
	</Alert>
{/if}
