<script lang="ts" context="module">
	import { Toggle } from 'flowbite-svelte';
	import type { ComponentProps } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { LabelOption } from './types';

	export type FormToggleOptions = Partial<HTMLInputAttributes & ComponentProps<Toggle>> &
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
	export let options: FormToggleOptions = {};

	const { value, errors, constraints } = formFieldProxy(superform, field as string);
</script>

<div class="space-y-2">
	<div class="flex items-center">
		<Toggle
			{...options}
			bind:checked={$value}
			name={field}
			value="true"
			data-invalid={$errors}
			{...$constraints}
		/>
		<Label color={$errors ? 'red' : 'gray'}>
			<div>
				<span><slot>{options.label}</slot></span>
				<FieldRequiredIndicator {field} />
			</div>
		</Label>
	</div>
	<FieldError {field} />
</div>
