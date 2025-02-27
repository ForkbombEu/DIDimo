<script lang="ts">
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field } from '@/forms/fields';
	import { goto, m } from '@/i18n';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { pb } from '@/pocketbase';
	import type { CredentialIssuersRecord, Data, ServicesRecord } from '@/pocketbase/types';
	import { nanoid } from 'nanoid';

	//

	const form = createForm({
		adapter: zod(
			z.object({
				url: z.string().url()
			})
		),
		onSubmit: async ({ form }) => {
			// note: use npm:out-of-character to clean the url if needed
			const { url } = form.data;

			const credentialIssuer = await pb
				.collection('credential_issuers')
				.create({ url } satisfies Data<CredentialIssuersRecord>);

			const service = await pb.collection('services').create({
				name: nanoid(5),
				credential_issuers: [credentialIssuer.id]
			} satisfies Data<ServicesRecord>);

			goto(`/services/${service.id}`);
		}
	});
</script>

<Form {form} hideRequiredIndicator>
	<Field {form} name="url" options={{ type: 'url', label: m.Credential_issuer_URL() }} />

	{#snippet submitButton()}
		<SubmitButton class="flex w-full">{m.Start_check()}</SubmitButton>
	{/snippet}
</Form>
