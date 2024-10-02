<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { goto } from '$lib/i18n';
	import { pb } from '$lib/pocketbase';
	import { Collections } from '$lib/pocketbase/types';
	import { Form, createForm, FormError, SubmitButton, Input } from '$lib/forms';
	import { z } from 'zod';
	import { currentEmail } from './+layout.svelte';

	const schema = z.object({
		email: z.string().email(),
		password: z.string()
	});

	const superform = createForm(
		schema,
		async ({ form }) => {
			const { data } = form;
			const u = pb.collection('users');
			await u.authWithPassword(data.email, data.password);
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

	<Input
		{superform}
		field="password"
		options={{
			id: 'password',
			type: 'password',
			label: 'Your password',
			placeholder: '•••••'
		}}
	/>

	<FormError />

	<div class="flex justify-end">
		<SubmitButton>Log in</SubmitButton>
	</div>
</Form>
