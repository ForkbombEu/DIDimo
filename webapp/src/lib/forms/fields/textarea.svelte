<script lang="ts" context="module">
	import { Textarea } from 'flowbite-svelte';
	import type { ComponentProps } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { LabelOption } from './types';

	export type FormTextareaOptions = Partial<HTMLInputAttributes & ComponentProps<Textarea>> &
		LabelOption;
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';
	import FieldWrapper from './fieldParts/fieldWrapper.svelte';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: FormTextareaOptions = {};

	const { value, errors, constraints } = formFieldProxy(superform, field as string);
</script>

<FieldWrapper {field} label={options.label}>
	<Textarea
		{...options}
		bind:value={$value}
		class={`!min-h-[200px] ${options.class}`}
		color={$errors ? 'red' : 'base'}
		name={field}
		data-invalid={$errors}
		{...$constraints}
	/>
</FieldWrapper>
