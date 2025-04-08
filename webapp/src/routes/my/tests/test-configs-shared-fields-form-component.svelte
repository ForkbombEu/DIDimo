<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Form } from '@/forms';
	import FieldConfigToFormField from './field-config-to-form-field.svelte';
	import { TestConfigsSharedFieldsForm } from './test-configs-shared-fields-form.svelte.js';

	//

	type Props = {
		form: TestConfigsSharedFieldsForm;
	};

	const { form }: Props = $props();
</script>

<Form form={form.form} hide={['submit_button']} hideRequiredIndicator>
	{#each form.fields as config}
		<FieldConfigToFormField {config} form={form.form} />
	{/each}
</Form>
