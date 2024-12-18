<script lang="ts">
	import { run } from 'svelte/legacy';

	import { goto, m } from '@/i18n';
	import { loginUser } from '@/webauthn/index';
	import z from 'zod';

	import { Form, createForm, FormError, SubmitButton } from '@/forms';
	import { Field } from '@/forms/fields';
	import { currentEmail } from '../+layout.svelte';
	import { zod } from 'sveltekit-superforms/adapters';

	const schema = z.object({
		email: z.string().email()
	});

	const form = createForm({
		adapter: zod(schema),
		onSubmit: async ({ form }) => {
			const { data } = form;
			await loginUser(data.email);
			await goto('/my');
		},
		initialData: { email: $currentEmail },
		options: { taintedMessage: null }
	});

	const { form: formData } = form;

	run(() => {
		$currentEmail = $formData.email;
	});
</script>

<Form {form}>
	<Field
		{form}
		name="email"
		options={{
			id: 'email',
			type: 'email',
			label: m.Your_email(),
			placeholder: 'name@foundation.org'
		}}
	/>

	{#snippet submitButtonContent()}
		{m.Log_in_with_webauthn()}
	{/snippet}
</Form>
