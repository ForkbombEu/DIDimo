<script lang="ts">
	import { CodeEditorField, Field } from '@/forms/fields';
	import type { FieldConfig } from './logic';
	import type { SuperForm } from 'sveltekit-superforms';

	//

	type Props = {
		config: FieldConfig;
		form: SuperForm<Record<string, unknown>>;
	};

	const { config, form }: Props = $props();
</script>

{#if config.Type == 'string'}
	<Field {form} name={config.CredimiID} options={{ label: config.LabelKey }} />
{:else if config.Type == 'object'}
	<CodeEditorField
		{form}
		name={config.CredimiID}
		options={{
			lang: 'json',
			label: config.LabelKey,
			value: config.Example
		}}
	/>
{/if}
