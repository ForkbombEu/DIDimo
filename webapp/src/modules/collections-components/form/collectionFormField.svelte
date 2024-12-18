<script lang="ts" module>
	import type { CollectionName, AnySchemaField } from '@/pocketbase/collections-models';
	import type { FieldSnippet, RelationFieldOptions } from './collectionFormTypes';

	export type CollectionFormFieldProps<C extends CollectionName> = {
		fieldConfig: AnySchemaField;
		hidden?: boolean;
		label?: string;
		snippet?: FieldSnippet<C>;
		relationFieldOptions?: RelationFieldOptions<C>;
		description?: string;
		placeholder?: string;
	};
</script>

<script lang="ts" generics="C extends CollectionName">
	import { getFormContext } from '@/forms';
	import { CheckboxField, FileField, Field, SelectField, TextareaField } from '@/forms/fields';
	import CollectionField from '../collectionField.svelte';
	import { getCollectionNameFromId } from '@/pocketbase/collections-models';
	import { isArrayField } from '@/pocketbase/collections-models';
	import type { FormPath, SuperForm } from 'sveltekit-superforms';
	import type { CollectionFormData } from '@/pocketbase/types';

	//

	let {
		fieldConfig,
		label = fieldConfig.name,
		description,
		placeholder,
		hidden = false,
		snippet,
		relationFieldOptions = {}
	}: CollectionFormFieldProps<C> = $props();

	//

	const name = fieldConfig.name;
	const multiple = isArrayField(fieldConfig);
	const { form } = getFormContext();
</script>

{#if hidden}
	<!-- Nothing -->
{:else if snippet}
	{@render snippet({
		form: form as SuperForm<CollectionFormData[C]>,
		field: name as FormPath<CollectionFormData[C]>
	})}
{:else if fieldConfig.type == 'text' || fieldConfig.type == 'url' || fieldConfig.type == 'date' || fieldConfig.type == 'email'}
	<Field {form} {name} options={{ label, description, placeholder, type: fieldConfig.type }} />
{:else if fieldConfig.type == 'number'}
	<Field {form} {name} options={{ label, description, type: 'number', placeholder }} />
{:else if fieldConfig.type == 'json'}
	<TextareaField {form} {name} options={{ label, description, placeholder }} />
{:else if fieldConfig.type == 'bool'}
	<CheckboxField {form} {name} options={{ label, description }} />
{:else if fieldConfig.type == 'file'}
	{@const accept = fieldConfig.options.mimeTypes?.join(',')}
	<FileField {form} {name} options={{ label, multiple, accept, placeholder }} />
{:else if fieldConfig.type == 'select'}
	{@const items = fieldConfig.options.values?.map((v) => ({ label: v, value: v }))}
	<SelectField
		{form}
		{name}
		options={{ label, items, type: multiple ? 'multiple' : 'single', description, placeholder }}
	/>
{:else if fieldConfig.type == 'editor'}
	<TextareaField {form} {name} options={{ label, description, placeholder }} />
{:else if fieldConfig.type == 'relation'}
	{@const collectionName = getCollectionNameFromId(fieldConfig.options.collectionId!) as C}
	<CollectionField
		{form}
		{name}
		collection={collectionName}
		options={{
			...relationFieldOptions,
			multiple,
			label,
			description
		}}
	/>
{/if}
