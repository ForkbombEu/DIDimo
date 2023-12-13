<script lang="ts">
	import { Button, Heading, P } from 'flowbite-svelte';
	import { currentUser } from '$lib/pocketbase';
	import CopyButton from '$lib/components/copyButton.svelte';
	import { capitalizeFirstLetter } from '$lib/utils/strings';
	import type { UsersRecord } from '$lib/pocketbase/types';

	type UserKeys = keyof UsersRecord;

	const keys: Array<UserKeys> = [
		'ethereum_address',
		'ecdh_public_key',
		'eddsa_public_key',
		'reflow_public_key',
		'bitcoin_public_key'
	];
	const hasKeys = keys.map((k) => $currentUser?.[k]).every((k) => Boolean(k));
</script>

<Heading tag="h6" class="mb-1">Your keys</Heading>
{#if hasKeys}
	<div class="flex flex-col gap-4">
		{#each keys as key}
			{@const value = $currentUser?.[key]}
			{#if value}
				{@const title = capitalizeFirstLetter(key.replaceAll('_', ' '))}
				<div class="flex flex-row items-center justify-between w-full gap-4">
					<div class="w-0 grow overflow-hidden">
						<div class="text-md text-black font-semibold max-w-sm">{title}</div>
						<div class="text-ellipsis overflow-hidden">{value}</div>
					</div>
					<div class="shrink-0">
						<CopyButton textToCopy={value} />
					</div>
				</div>
			{/if}
		{/each}
	</div>
{:else}
	<P>You have no keys yet.</P>
	<div class="flex justify-end mt-4">
		<Button color="primary" size="sm" href="/keypairoom/regenerate">Generate new keys</Button>
	</div>
{/if}
