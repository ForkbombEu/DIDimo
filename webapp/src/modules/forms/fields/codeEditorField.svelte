<script lang="ts" generics="Data extends GenericRecord">
	import type { GenericRecord } from '@/utils/types';
	import * as Form from '@/components/ui/form';
	import type { FormPathLeaves, SuperForm } from 'sveltekit-superforms';
	import { stringProxy } from 'sveltekit-superforms';
	import type { ComponentProps } from 'svelte';
	import FieldWrapper from './parts/fieldWrapper.svelte';
	import type { FieldOptions } from './types';
	import CodeEditor from '@/components/ui-custom/codeEditor.svelte';

	//

	interface Props {
		form: SuperForm<Data>;
		name: FormPathLeaves<Data, string | number>;
		options: Partial<FieldOptions> & ComponentProps<typeof CodeEditor>;
	}

	const { form, name, options }: Props = $props();

	const { form: formData } = $derived(form);
	const valueProxy = $derived(stringProxy(formData, name, { empty: 'undefined' }));
</script>

<Form.Field {form} {name}>
	<FieldWrapper field={name} {options}>
		{#snippet children({ props })}
			<CodeEditor {...options} {...props} bind:value={$valueProxy} />
		{/snippet}
	</FieldWrapper>
</Form.Field>
