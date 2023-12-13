<script lang="ts">
	import { Form, createForm } from '$lib/forms';
	import FieldController from '$lib/forms/fields/fieldController.svelte';
	import { z } from 'zod';
	import { Input, Button, Heading } from 'flowbite-svelte';
	import Card from '$lib/components/card.svelte';

	const schema = z.object({ url: z.string().url() });

	const superform = createForm(schema, handleSubmit);

	async function handleSubmit() {}
</script>

<Card class="p-8">
	<Form {superform} className="space-y-8">
		<Heading>Validate your Credential Issuer!</Heading>
		<FieldController {superform} field="url" let:value let:updateValue>
			<Input
				{value}
				size="lg"
				on:change={(e) => updateValue(e.target.value)}
				placeholder="Enter your credential issuer url"
			/>
		</FieldController>
		<Button size="xl" class="w-full">Check</Button>
	</Form>
</Card>
