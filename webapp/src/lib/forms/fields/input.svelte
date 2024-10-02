<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	import { Input } from 'flowbite-svelte';
	import type { HelpTextOption, LabelOption } from './types';
	import type { ComponentProps } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	export type FormInputOptions = Partial<HTMLInputAttributes & ComponentProps<Input>> &
		LabelOption &
		HelpTextOption;
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';
	import FieldWrapper from './fieldParts/fieldWrapper.svelte';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: FormInputOptions = {};

	let type = options.type ?? 'text';

	const { value, errors, constraints } = formFieldProxy(superform, field as string);
</script>

<FieldWrapper {field} label={options.label} helpText={options.helpText}>
	<slot slot="labelRight" name="labelRight" />
	<Input
		{...options}
		{type}
		color={$errors ? 'red' : 'base'}
		name={field}
		data-invalid={$errors}
		bind:value={$value}
		{...$constraints}
	/>
</FieldWrapper>
