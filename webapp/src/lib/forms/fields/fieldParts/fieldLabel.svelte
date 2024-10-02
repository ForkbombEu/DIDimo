<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { getFormContext } from '../../form.svelte';
	import { formFieldProxy } from 'sveltekit-superforms/client';
	import { fieldHasErrors } from './fieldError.svelte';
	import { Label } from 'flowbite-svelte';
	import FieldRequiredIndicator from './fieldRequiredIndicator.svelte';
	import { capitalizeFirstLetter } from '$lib/utils/strings';

	export let field: string;
	export let text: string | null | undefined = undefined;
	if (!text) text = capitalizeFirstLetter(field);

	const { superform } = getFormContext();
	const { errors } = formFieldProxy(superform, field);
</script>

<Label
	for={field}
	class="label flex items-center justify-between"
	color={fieldHasErrors($errors) ? 'red' : 'gray'}
>
	<div>
		<span>{text}</span>
		<FieldRequiredIndicator {field} />
	</div>

	<slot name="right" />
</Label>
