<script lang="ts">
	import { CodeEditorField } from '@/forms/fields';
	import { createForm, Form } from '@/forms';
	import { zod } from 'sveltekit-superforms/adapters';
	import {
		createSchemaFromFieldsConfigs,
		stringifiedObjectSchema,
		type SpecificFieldConfig,
		type TestInput
	} from './logic';
	import { Record as R, Record, pipe } from 'effect';
	import { Store, watch } from 'runed';
	import FieldConfigToFormField from './field-config-to-form-field.svelte';
	import { Input } from '@/components/ui/input';
	import Label from '@/components/ui/label/label.svelte';
	import IconButton from '@/components/ui-custom/iconButton.svelte';
	import { Pencil, Lock, Info } from 'lucide-svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import Alert from '@/components/ui-custom/alert.svelte';
	import { Button } from '@/components/ui/button';
	import Separator from '@/components/ui/separator/separator.svelte';
	import { nanoid } from 'nanoid';

	//

	type UpdateFunction = (testInput: TestInput) => void;

	type Props = {
		fields: SpecificFieldConfig[];
		jsonConfig: Record<string, unknown>;
		defaultFieldsIds?: string[];
		defaultValues?: Record<string, unknown>;
		onValidUpdate?: UpdateFunction;
	};

	const {
		fields,
		defaultValues = {},
		defaultFieldsIds = [],
		onValidUpdate,
		jsonConfig = {}
	}: Props = $props();

	console.log(defaultFieldsIds);

	/* Form creation */

	const JSON_CONFIG_KEY = 'jsonConfig';
	const jsonConfigString = JSON.stringify(jsonConfig, null, 4);

	const form = createForm({
		adapter: zod(
			createSchemaFromFieldsConfigs(fields).extend({
				jsonConfig: stringifiedObjectSchema.optional()
			})
		),
		initialData: {
			[JSON_CONFIG_KEY]: jsonConfigString
		},
		options: {
			id: nanoid(6)
		}
	});

	const { form: formData } = form;

	/* Overrides */

	let overrideFields = $state<string[]>([]);

	function resetOverride(id: string) {
		overrideFields = overrideFields.filter((f) => f !== id);
		$formData = { ...$formData, [id]: defaultValues[id] };
	}

	/* Update form data when default values change */

	watch(
		() => defaultValues,
		() => {
			const notOverridden = pipe(
				defaultValues,
				// Only fields that are in the fields array
				R.filter((_, key) => fields.map((f) => f.CredimiID).includes(key)),
				// Not overridden
				R.filter((_, key) => !overrideFields.includes(key))
			);
			$formData = { ...$formData, ...notOverridden };
		}
	);

	/* Disable fields when jsonConfig is manually edited */

	const { tainted } = form;

	const taintedState = new Store(tainted);
	const isJsonConfigTainted = $derived(Boolean(taintedState.current?.jsonConfig));

	function resetJsonConfig() {
		$formData = { ...$formData, jsonConfig: jsonConfigString };
	}

	/* Trigger onValidUpdate */

	const { validateForm, validate } = form;
	const formState = new Store(formData); // Readonly

	$effect(() => {
		if (isJsonConfigTainted) {
			const jsonConfigString = formState.current[JSON_CONFIG_KEY] as string;

			validate(JSON_CONFIG_KEY, { update: false, value: jsonConfigString }).then((errors) => {
				if (Boolean(errors?.length)) return;
				else
					onValidUpdate?.({
						format: 'json',
						data: JSON.parse(jsonConfigString)
					});
			});
		} else {
			const testInput: TestInput = {
				format: 'variables',
				data: R.remove(formState.current, JSON_CONFIG_KEY)
			};
			validateForm({ update: false }).then((result) => {
				if (result.valid) onValidUpdate?.(testInput);
			});
		}
	});
</script>

<Form
	{form}
	hide={['submit_button']}
	hideRequiredIndicator
	class="flex flex-col gap-16 md:flex-row"
>
	<div class="shrink-0 grow basis-1">
		<CodeEditorField
			{form}
			name={JSON_CONFIG_KEY}
			options={{
				lang: 'json',
				label: 'JSON configuration'
			}}
		/>
	</div>

	<div class="shrink-0 grow basis-1">
		<div class="mb-8 space-y-2">
			<Label>Fields</Label>
			<Separator />
		</div>
		{#if isJsonConfigTainted}
			<div class="text-muted-foreground text-sm">
				<Alert variant="info" icon={Info}>
					{#snippet content({ Title, Description })}
						<Title class="font-bold">Info</Title>
						<Description class="mb-2">
							JSON configuration is manually edited. Fields are disabled.
						</Description>

						<Button
							variant="outline"
							onclick={() => {
								resetJsonConfig();
							}}
						>
							Reset JSON and use fields
						</Button>
					{/snippet}
				</Alert>
			</div>
		{:else}
			<div class="space-y-8">
				{#each fields as config}
					{@const { CredimiID: id, LabelKey: label } = config}
					{@const isDefault = defaultFieldsIds.includes(id)}
					{@const isOverride = overrideFields.includes(id)}
					<pre>{config.CredimiID} - {config.LabelKey} - {config.FieldName}</pre>

					{#if isDefault && !isOverride}
						<div class="space-y-2">
							<Label>{label}</Label>
							<div class="flex items-center gap-2">
								<Input
									value={defaultValues[id]}
									disabled
									placeholder="Value depends on shared '{label}' field"
									class={{
										'font-mono':
											config.Type == 'object' && Boolean(defaultValues[id])
									}}
								/>
								<IconButton
									icon={Pencil}
									variant="outline"
									onclick={() => {
										overrideFields.push(id);
									}}
								/>
							</div>
						</div>
					{:else}
						<div class="relative">
							{#if isOverride}
								<div class="absolute right-0 top-0">
									<button
										class="text-primary flex items-center gap-2 text-sm underline hover:no-underline"
										onclick={() => {
											resetOverride(id);
										}}
									>
										<Icon src={Lock} size={14} />
										Lock again
									</button>
								</div>
							{/if}
							<FieldConfigToFormField {config} {form} />
						</div>
					{/if}
				{/each}
			</div>
		{/if}
	</div>
</Form>
