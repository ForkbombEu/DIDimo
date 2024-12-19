<script lang="ts">
	import { goto, m } from '@/i18n';
	import { pb } from '@/pocketbase';
	import { Form, createForm } from '@/forms';
	import { Field } from '@/forms/fields';
	import z from 'zod';
	import { currentEmail } from './+layout.svelte';
	import { zod } from 'sveltekit-superforms/adapters';

	//

	const schema = z.object({
		email: z.string().email(),
		password: z.string()
	});

	const form = createForm({
		adapter: zod(schema),
		onSubmit: async ({ form }) => {
			const { data } = form;
			const u = pb.collection('users');
			await u.authWithPassword(data.email, data.password);
			await goto('/my');
		},
		initialData: { email: currentEmail.value },
		options: { taintedMessage: null }
	});

	const { form: formData } = form;

	$effect(() => {
		currentEmail.value = $formData.email;
	});
</script>

<Form {form}>
	<Field
		{form}
		name="email"
		options={{
			type: 'email',
			label: m.Your_email(),
			placeholder: 'name@foundation.org'
		}}
	/>

	<Field
		{form}
		name="password"
		options={{
			type: 'password',
			label: m.Your_password(),
			placeholder: '•••••'
		}}
	/>

	{#snippet submitButtonContent()}
		{m.Log_in()}
	{/snippet}
</Form>
