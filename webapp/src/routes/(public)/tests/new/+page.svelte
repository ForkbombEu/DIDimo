<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field } from '@/forms/fields';
	import { m } from '@/i18n';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';

	//

	const form = createForm({
		adapter: zod(
			z.object({
				url: z.string().url()
			})
		),
		onSubmit: async ({ form }) => {
			const { url } = form.data;
			// Run request here
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
