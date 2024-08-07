<script lang="ts">
	import { Button, Heading, P } from 'flowbite-svelte';
	import CopyButton from '$lib/components/copyButton.svelte';
	import { capitalizeFirstLetter } from '$lib/utils/strings';
	import type { PublicKeys } from '$lib/keypairoom/utils';

	export let keys: PublicKeys | undefined = undefined;
</script>

<Heading tag="h6" class="mb-1">Your keys</Heading>
{#if keys}
	<div class="flex flex-col gap-4">
		{#each Object.entries(keys) as [keyName, key]}
			{@const title = capitalizeFirstLetter(keyName.replaceAll('_', ' '))}
			<div class="flex flex-row items-center justify-between w-full gap-4">
				<div class="w-0 grow overflow-hidden">
					<div class="text-md text-black font-semibold max-w-sm">{title}</div>
					<div class="text-ellipsis overflow-hidden">{key}</div>
				</div>
				<div class="shrink-0">
					<CopyButton textToCopy={key} />
				</div>
			</div>
		{/each}
	</div>
{:else}
	<P>You have no keys yet.</P>
	<div class="flex justify-end mt-4">
		<Button color="primary" size="sm" href="/keypairoom">Generate keyring</Button>
	</div>
{/if}
