<script lang="ts">
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field, TextareaField } from '@/forms/fields';
	import { m } from '@/i18n';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { pb } from '@/pocketbase';
	import { CollectionField } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';

	//

	let result = $state<string>();

	const form = createForm({
		adapter: zod(
			z.object({
				name: z.string(),
				standard: z.string(),
				json: z.string(),
				email: z.string().email()
			})
		),
		onSubmit: async ({ form }) => {
			const { name, standard, json, email } = form.data;

			// /pkg/internal/pb/workflow.go
			const { message } = await pb.send('/api/openid4vp-test', {
				body: {
					input: json,
					user_mail: email,
					workflow_id: standard,
					test_name: name
				}
			});

			result = message;
		}
	});
</script>

<Form {form} hideRequiredIndicator>
	<Field {form} name="name" options={{ type: 'text', label: m.Wallet_name() }} />
	<CollectionField {form} name="standard" collection="standards" />
	<TextareaField {form} name="json"></TextareaField>

	<div class="bg-secondary/30 space-y-4 rounded-xl border p-4">
		<T>{m.Send_the_results_to_this_email()}</T>
		<Field {form} name="email" options={{ type: 'email' }} />
	</div>

	{#snippet submitButton()}
		<SubmitButton class="flex w-full">{m.Start_check()}</SubmitButton>
	{/snippet}

	{#if result}
		{result}
	{/if}
</Form>
