<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { goto } from '$lib/i18n';
	import { loginUser } from '$lib/webauthn/index';
	import { z } from 'zod';

	import { Form, createForm, Input, FormError, SubmitButton } from '$lib/forms';
	import { currentEmail } from '../+layout.svelte';

	const schema = z.object({
		email: z.string().email()
	});

	const superform = createForm(
		schema,
		async ({ form }) => {
			const { data } = form;
			await loginUser(data.email);
			await goto('/my');
		},
		{ email: $currentEmail },
		{ taintedMessage: null }
	);

	const { capture, restore, form } = superform;
	export const snapshot = { capture, restore };

	$: $currentEmail = $form.email;
</script>

<Form {superform}>
	<Input
		{superform}
		field="email"
		options={{
			id: 'email',
			type: 'email',
			label: 'Your email',
			placeholder: 'name@foundation.org'
		}}
	/>

	<FormError />

	<div class="flex justify-end">
		<SubmitButton>Log in with webauthn</SubmitButton>
	</div>
</Form>
