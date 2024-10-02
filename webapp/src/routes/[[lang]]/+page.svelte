<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import FieldController from '$lib/forms/fields/fieldController.svelte';
	import { Input, Button, Heading, Spinner, Alert, P } from 'flowbite-svelte';
	import Card from '$lib/components/card.svelte';
	import { superForm } from 'sveltekit-superforms/client';
	import FeaturesDisplay from './_lib/FeaturesDisplay.svelte';
	import { InformationCircle } from 'svelte-heros-v2';

	//

	export let data;
	export let form;

	const superform = superForm(data.form, {
		taintedMessage: null
	});
	const { enhance, delayed } = superform;
</script>

<Card class="p-8 ">
	<form use:enhance method="post" class="space-y-10">
		<Heading>Validate your Credential Issuer!</Heading>
		{#if !$delayed}
			<FieldController {superform} field="url" let:value let:updateValue>
				<Input
					name="url"
					size="lg"
					placeholder="Enter your credential issuer URL here!"
					{value}
					on:change={(e) => updateValue(e?.currentTarget?.value)}
				/>
			</FieldController>

			{#if form?.connectionError}
				<Alert dismissable color="red" border>
					<svelte:fragment slot="icon">
						<InformationCircle size="20" />
					</svelte:fragment>
					<span class="font-semibold mr-2">Connection error!</span>
				</Alert>
			{/if}

			{#if form?.error}
				<Alert dismissable color="red" border>
					<svelte:fragment slot="icon">
						<InformationCircle size="20" />
					</svelte:fragment>
					<span class="font-semibold mr-2">{form.error}</span>
				</Alert>
			{/if}

			{#if form?.features}
				<FeaturesDisplay features={form.features} />
			{/if}

			{#if !form?.features}
				<Button size="xl" class="w-full" type="submit">Check!</Button>
			{:else}
				<Button color="alternative" class="w-full" href="/">Restart!</Button>
			{/if}
		{:else if $delayed}
			<div class="flex flex-col items-center justify-center p-8">
				<Spinner />
				<P>Checking your stuff...</P>
			</div>
		{/if}
	</form>
</Card>
