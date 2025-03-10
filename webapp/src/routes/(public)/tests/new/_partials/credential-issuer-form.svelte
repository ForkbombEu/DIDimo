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

			const response = await pb.send('URL_HERE', {
				method: 'POST',
				body: {
					credentialIssuerUrl: url
				}
			});
			console.log(response);

			// goto(`/providers/${service.id}`);
		}
	});
</script>

<Form {form} hideRequiredIndicator>
	<Field {form} name="url" options={{ type: 'url', label: m.Credential_issuer_URL() }} />

	{#snippet submitButton()}
		<SubmitButton class="flex w-full">{m.Start_check()}</SubmitButton>
	{/snippet}
</Form>
