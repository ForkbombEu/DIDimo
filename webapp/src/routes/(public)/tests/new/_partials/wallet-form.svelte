<script lang="ts">
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field, TextareaField } from '@/forms/fields';
	import { m } from '@/i18n';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { pb } from '@/pocketbase';
	import { CollectionField } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';

	import t from './test.json';

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
		initialData: {
			name: 'Wallet test',
			standard: 'w742h03qg7c606i',
			email: 'pin@gmail.com',
			json: JSON.stringify(t, null, 4)
		},
		onSubmit: async ({ form }) => {
			try {
				const { name, standard, json, email } = form.data;

				// /pkg/internal/pb/workflow.go
				const result = await pb.send('/api/openid4vp-test', {
					method: 'POST',
					body: {
						input: JSON.parse(json),
						user_mail: email,
						workflow_id: standard,
						tesst_name: name
					}
				});
				console.log(result);

				// result = message;
			} catch (e) {
				console.log(e);
			}
		}
	});
</script>

<Form {form} hideRequiredIndicator>
	<Field {form} name="name" options={{ type: 'text', label: m.Wallet_name() }} />
	<CollectionField
		{form}
		name="standard"
		collection="standards"
		options={{ displayFields: ['name'] }}
	/>
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
