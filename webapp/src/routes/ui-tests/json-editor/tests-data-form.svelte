<script lang="ts">
	// import { sharedFields, testsFields, testsConfigs } from './sample.data';
	import FieldConfigFormShared from './field-config-form-shared.svelte';
	import FieldConfigForm from './field-config-form.svelte';
	import type { FieldsResponse, TestInput } from './logic';

	//

	type Props = {
		data: FieldsResponse;
	};

	let { data }: Props = $props();

	//

	let sharedData = $state<Record<string, unknown>>({});

	const defaultFieldsIds = Object.values(data.normalized_fields).map((f) => f.CredimiID);
	console.log(defaultFieldsIds);

	const masterDataStructure: Record<string, TestInput> = $state({});
</script>

<div class="space-y-16">
	<div class="space-y-4">
		<h2 class="text-lg font-bold">Shared fields</h2>
		<FieldConfigFormShared
			fields={data.normalized_fields}
			onUpdate={(form) => (sharedData = form)}
		/>
	</div>

	<hr />
	{#each Object.entries(data.specific_fields) as [testId, testData]}
		<div class="space-y-4">
			<h2 class="text-lg font-bold">{testId}</h2>
			<FieldConfigForm
				fields={testData.fields}
				jsonConfig={JSON.parse(testData.content)}
				defaultValues={sharedData}
				{defaultFieldsIds}
				onValidUpdate={(form) => {
					masterDataStructure[testId] = form;
				}}
			/>
		</div>
		<hr />
	{/each}
</div>
