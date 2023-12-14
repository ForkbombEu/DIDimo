<script lang="ts">
	import FieldController from '$lib/forms/fields/fieldController.svelte';
	import { Input, Button, Heading, Spinner, Alert, P } from 'flowbite-svelte';
	import Card from '$lib/components/card.svelte';
	import { superForm } from 'sveltekit-superforms/client';
	import { InformationCircle } from 'svelte-heros-v2';
	import { CredentialIssuersScansErrorOptions as Errors } from '$lib/pocketbase/types.js';

	//

	export let data;
	export let form;

	const superform = superForm(data.form, {
		taintedMessage: null
	});
	const { enhance, delayed, form: f } = superform;

	$f.url = 'https://raw.githubusercontent.com/ForkbombEu/DIDroom_microservices/main/public/';
	// $f.url = 'http://0.0.0.0:3000';

	const errorMessages: Record<Errors, string> = {
		[Errors.BAD_JSON]: 'JSON badly formatted.',
		[Errors.CONNECTION_ERROR]: "Couldn't manage to connect to the url.",
		[Errors.VALIDATION_ERROR]: 'The metadata file is not compliant to the specification.',
		[Errors.FILE_NOT_FOUND]: 'The metadata file does not exist at the required location.'
	};
</script>

<Card class="p-8 space-y-8">
	<Heading>Validate your Credential Issuer!</Heading>
	<form use:enhance method="post" class="space-y-6">
		{#if !$delayed}
			<FieldController {superform} field="url" let:value let:updateValue>
				<Input
					name="url"
					size="lg"
					placeholder="Enter your credential issuer URL here!"
					{value}
					on:change={(e) => updateValue(e.target.value)}
				/>
			</FieldController>

			{#if form?.error}
				<Alert dismissable color="red" border>
					<svelte:fragment slot="icon">
						<InformationCircle size="20" />
					</svelte:fragment>
					<span class="font-semibold mr-2">Error!</span>
					{errorMessages[form.error] ?? form.error}
				</Alert>
			{/if}

			{#if form?.success}
				<Alert color="green" border>
					<svelte:fragment slot="icon">
						<InformationCircle size="20" />
					</svelte:fragment>
					<span class="font-semibold mr-2">All good!</span>
				</Alert>
			{/if}

			<Button size="xl" class="w-full" type="submit">Check!</Button>
		{:else if $delayed}
			<Spinner></Spinner>
			<P>Checking your stuff...</P>
		{/if}
	</form>
</Card>
