<script lang="ts">
	import Input from '@/components/ui/input/input.svelte';
	import { Form, type FieldData, type Field } from './logic.exp.svelte';
	import { z } from 'zod';

	const sharedFields: FieldData[] = [
		{ id: 'name', schema: z.string(), label: 'Name', description: 'The name of the user' },
		{ id: 'age', schema: z.number().min(18), label: 'Age', description: 'The age of the user' }
	];

	const primaryForm = new Form({
		fieldsData: sharedFields,
		defaultValues: () => ({}),
		onSubmit: () => {}
	});

	const form = new Form({
		fieldsData: [
			...sharedFields,
			{
				id: 'email',
				schema: z.string().email(),
				label: 'Email',
				description: 'The email of the user'
			},
			{
				id: 'password',
				schema: z.string().min(8),
				label: 'Password',
				description: 'The password of the user'
			}
		],
		defaultValues: () => primaryForm.validData,
		onSubmit: () => {}
	});
</script>

<div class="flex flex-col gap-16">
	<div>
		<h1>Primary Form</h1>
		<pre>{JSON.stringify(primaryForm.validData, null, 2)}</pre>

		<form>
			{#each primaryForm.fields as field}
				{@render fieldRenderer(field)}
			{/each}
		</form>
	</div>

	<hr />

	<div>
		<h1>Form</h1>
		<pre>{JSON.stringify(form.validData, null, 2)}</pre>

		<form>
			{#each form.fields as field}
				{@render fieldRenderer(field)}
			{/each}
		</form>
	</div>
</div>

<!--  -->
{#snippet fieldRenderer(field: Field)}
	{@const inputType = field.getSchemaType()}
	<div class="flex flex-col gap-2">
		<label for={field.id}>{field.label}</label>
		{#if inputType == 'text'}
			<Input bind:value={field.value} />
		{:else if inputType == 'number'}
			<Input bind:value={field.value} type="number" />
		{:else}
			<p>not defined</p>
		{/if}
		{#if field.error}
			<p class="text-red-500">{field.error}</p>
		{/if}
	</div>
{/snippet}
