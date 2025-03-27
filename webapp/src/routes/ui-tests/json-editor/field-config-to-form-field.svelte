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

{#if config.field_type == 'string'}
	<Field {form} name={config.credimi_id} options={{ label: config.i18_label }} />
{:else if config.field_type == 'object'}
	<CodeEditorField
		{form}
		name={config.credimi_id}
		options={{
			lang: 'json',
			label: config.i18_label
		}}
	/>
{/if}
