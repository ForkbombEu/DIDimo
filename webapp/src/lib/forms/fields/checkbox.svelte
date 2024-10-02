<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	import { Checkbox } from 'flowbite-svelte';
	import type { ComponentProps } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { LabelOption } from './types';

	export type FormCheckboxOptions = Partial<HTMLInputAttributes & ComponentProps<Checkbox>> &
		LabelOption;
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	import { Label } from 'flowbite-svelte';
	import FieldError from './fieldParts/fieldError.svelte';
	import FieldRequiredIndicator from './fieldParts/fieldRequiredIndicator.svelte';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: FormCheckboxOptions = {};

	const { value, errors, constraints } = formFieldProxy(superform, field as string);
</script>

<div class="space-y-2">
	<Label color={$errors ? 'red' : 'gray'}>
		<div class="flex items-center space-x-2">
			<Checkbox
				{...options}
				color={$errors ? 'red' : 'secondary'}
				bind:checked={$value}
				name={field}
				value="true"
				data-invalid={$errors}
				{...$constraints}
			/>
			<div>
				<span><slot>{options.label}</slot></span>
				<FieldRequiredIndicator {field} />
			</div>
		</div>
	</Label>
	<FieldError {field} />
</div>
