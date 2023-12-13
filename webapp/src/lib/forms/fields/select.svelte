<script lang="ts" context="module">
	import type { LabelOption } from './types';

	import type { HTMLInputAttributes } from 'svelte/elements';

	export type FormSelectOptions = {
		[K in keyof HTMLInputAttributes]?: Exclude<HTMLInputAttributes[K], null>;
	} & {
		options?: string[];
		size?: 'sm' | 'lg' | 'md' | undefined;
	} & LabelOption;
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	import { Select, MultiSelect } from 'flowbite-svelte';
	import type { SelectOptionType } from 'flowbite-svelte/dist/types';
	import FieldWrapper from './fieldParts/fieldWrapper.svelte';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: FormSelectOptions = {};

	let { options: items = [], multiple = false, label = '' } = options;

	const { value, errors, constraints } = formFieldProxy(superform, field as string);

	const selectOptions: Array<SelectOptionType<string>> = items.map((i) => {
		return { value: i, name: i };
	});
</script>

<FieldWrapper {field} {label}>
	{#if !multiple}
		<Select
			{...options}
			name={field}
			items={selectOptions}
			bind:value={$value}
			color={$errors ? 'red' : 'base'}
			data-invalid={$errors}
			{...$constraints}
		/>
	{:else}
		<MultiSelect
			{...options}
			name={field}
			items={selectOptions}
			bind:value={$value}
			data-invalid={$errors}
			dropdownClass="bg-white dark:bg-gray-800"
			{...$constraints}
		/>
	{/if}
</FieldWrapper>
