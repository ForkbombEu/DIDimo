<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
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
			const { url } = form.data;

			const credentialIssuer = await pb
				.collection('credential_issuers')
				.create({ url } satisfies Data<CredentialIssuersRecord>);

			const provider = await pb.collection('services').create({
				name: nanoid(5),
				credential_issuers: [credentialIssuer.id]
			} satisfies Data<ServicesRecord>);

			goto(`/providers/${provider.id}`);
		}
	});
</script>

<div class="mx-auto max-w-xl px-8 py-12">
	<T tag="h1" class="mb-12">{m.Start_a_new_compliance_check()}</T>
	<Form {form} hideRequiredIndicator>
		<Field {form} name="url" options={{ type: 'url', label: m.Credential_issuer_URL() }} />

		{#snippet submitButton()}
			<SubmitButton class="flex w-full">{m.Start_check()}</SubmitButton>
		{/snippet}
	</Form>
</div>
