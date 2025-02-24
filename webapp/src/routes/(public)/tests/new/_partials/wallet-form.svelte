<script lang="ts">
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field, TextareaField } from '@/forms/fields';
	import { goto, m } from '@/i18n';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { pb } from '@/pocketbase';
	import type { CredentialIssuersRecord, Data, ServicesRecord } from '@/pocketbase/types';
	import { nanoid } from 'nanoid';
	import { CollectionField } from '@/collections-components';

	//

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
			const credentialIssuer = await pb.send('/api/openid4vp-test', {
				body: {
					// TODO
				}
			});
			// const service = await pb.collection('services').create({
			// 	name: nanoid(5),
			// 	credential_issuers: [credentialIssuer.id]
			// } satisfies Data<ServicesRecord>);
			// goto(`/services/${service.id}`);
		}
	});
</script>

<Form {form} hideRequiredIndicator>
	<Field {form} name="name" options={{ type: 'text', label: m.Wallet_name() }} />
	<CollectionField {form} name="standard" collection="standards" />
	<TextareaField {form} name="json"></TextareaField>

	<div class="">
		<Field {form} name="email" options={{ type: 'email' }} />
	</div>

	{#snippet submitButton()}
		<SubmitButton class="flex w-full">{m.Start_check()}</SubmitButton>
	{/snippet}
</Form>
