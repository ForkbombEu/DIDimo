<script lang="ts">
	import { CheckCircle, ExclamationCircle } from 'svelte-heros-v2';
	import { CredentialIssuersFeaturesTypeOptions as Feature } from '$lib/pocketbase/types.js';
	import { Alert, Heading, Li, List, P } from 'flowbite-svelte';

	export let features: Feature[] = [];

	const featuresList: Array<{ feature: Feature; successText: string; failText: string }> = [
		{
			feature: Feature.FILE_EXISTS,
			successText: 'Metadata file found',
			failText: 'Metadata file not found'
		},
		{
			feature: Feature.VALID_JSON,
			successText: 'Is a valid JSON',
			failText: 'Your JSON file is not valid'
		},
		{
			feature: Feature.SCHEMA_COMPLIANT,
			successText: 'Respects the OpenID standard',
			failText: 'The JSON file does not respect the OpenID standard'
		}
	];

	$: allGood = features.length === featuresList.length;
</script>

<Alert color="gray" border>
	<div class="flex flex-col items-center justify-center">
		<div>
			<Heading tag="h4" class="mb-4">Here's your result:</Heading>
			<List tag="ul" class="!space-y-2">
				{#each featuresList as featureItem}
					{@const success = features.includes(featureItem.feature)}
					{@const color = success ? 'text-green-500' : 'text-red-500'}
					{@const icon = success ? CheckCircle : ExclamationCircle}
					{@const text = success ? featureItem.successText : featureItem.failText}
					<Li icon>
						<svelte:component this={icon} class={`mr-2 ${color}`} />
						<P {color}>{text}</P>
					</Li>
				{/each}
			</List>
		</div>
		{#if allGood}
			<div class="-rotate-2 bg-white w-fit px-8 py-4 border rounded-lg mt-6">
				<Heading tag="h2" color="text-green-500">🎉 &nbsp; All good! &nbsp; 🎉</Heading>
			</div>
		{/if}
	</div>
</Alert>
