<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Form } from '@/forms';
	import { type ConfigFieldSpecific } from './logic';
	import FieldConfigToFormField from './field-config-to-form-field.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import { Pencil, Undo, Eye } from 'lucide-svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import * as Popover from '@/components/ui/popover/index.js';
	import { TestConfigFieldsForm } from './test-config-fields-form.svelte.js';

	//

	type Props = {
		form: TestConfigFieldsForm;
	};

	const { form }: Props = $props();

	//

	function previewValue(value: unknown, type: ConfigFieldSpecific['Type']): string {
		const NULL_VALUE = '<null>';
		if (!value) return NULL_VALUE;
		if (type === 'string') return value as string;
		if (type === 'object') return JSON.stringify(JSON.parse(value as string), null, 4);
		return NULL_VALUE;
	}
</script>

<Form
	form={form.form}
	hide={['submit_button']}
	hideRequiredIndicator
	class="flex flex-col gap-8 md:flex-row"
>
	{#if form.isValid}
		<p class="text-sm text-green-500">Form is valid</p>
	{:else}
		<p class="text-sm text-red-500">Form is invalid</p>
	{/if}

	<pre>{JSON.stringify(form.currentData, null, 4)}</pre>

	{#each form.specificFields as config}
		<FieldConfigToFormField {config} form={form.form} />
	{/each}

	{#each form.currentOverriddenFields as config}
		<div class="relative">
			<div class="absolute right-0 top-0">
				<button
					class="text-primary flex items-center gap-2 text-sm underline hover:no-underline"
					onclick={(e) => {
						e.preventDefault(); // Important to prevent form submission
						form.resetOverride(config.CredimiID);
					}}
				>
					<Icon src={Undo} size={14} />
					Reset to default
				</button>
			</div>
			<FieldConfigToFormField {config} form={form.form} />
		</div>
	{/each}

	{#if form.currentSharedFields}
		<div class="space-y-2">
			<Label>Default fields</Label>
			<ul class="space-y-1">
				{#each form.currentSharedFields as { CredimiID, LabelKey, Type }}
					{@const value = form.sharedData()[CredimiID]}
					{@const valuePreview = previewValue(value, Type)}

					<li class="flex items-center gap-2">
						<span class="font-mono text-sm">{LabelKey}</span>

						<Popover.Root>
							<Popover.Trigger
								class="rounded-md p-1 hover:cursor-pointer hover:bg-gray-200"
							>
								<Eye size={14} />
							</Popover.Trigger>
							<Popover.Content class="dark overflow-auto">
								<pre class="text-xs">{valuePreview}</pre>
							</Popover.Content>
						</Popover.Root>

						<button
							class="rounded-md p-1 hover:cursor-pointer hover:bg-gray-200"
							onclick={(e) => {
								e.preventDefault(); // Important to prevent form submission
								form.overrideField(CredimiID);
							}}
						>
							<Pencil size={14} />
						</button>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</Form>
