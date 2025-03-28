<script lang="ts">
	import FieldConfigFormShared from './field-config-form-shared.svelte';
	import FieldConfigForm from './field-config-form.svelte';
	import { createTestListInputSchema, type FieldsResponse } from './logic';
	import { createForm, Form, SubmitButton } from '@/forms';
	import { zod } from 'sveltekit-superforms/adapters';
	import { Store } from 'runed';
	import * as Popover from '@/components/ui/popover';
	import Button from '@/components/ui/button/button.svelte';

	//

	type Props = {
		data: FieldsResponse;
	};

	let { data }: Props = $props();

	//

	let sharedData = $state<Record<string, unknown>>({});

	const defaultFieldsIds = Object.values(data.normalized_fields).map((f) => f.CredimiID);

	//

	const form = createForm({
		adapter: zod(createTestListInputSchema(data)),
		onSubmit: ({ form }) => {
			console.log(form.data);
		}
	});

	const { form: formData, validateForm } = form;
	const formState = new Store(formData);

	const testsIds = $derived(Object.keys(data.specific_fields));

	const incompleteTestsIdsPromise = $derived.by(() => {
		formState.current;
		return validateForm().then((result) => testsIds.filter((test) => test in result.errors));
	});

	const completeTestsCount = $derived(
		incompleteTestsIdsPromise.then((tests) => testsIds.length - tests.length)
	);

	const completionStatusPromise = $derived(
		Promise.all([completeTestsCount, incompleteTestsIdsPromise])
	);

	//

	const SHARED_FIELDS_ID = 'shared-fields';
</script>

<div class="space-y-16">
	<div class="space-y-4">
		<h2 id={SHARED_FIELDS_ID} class="text-lg font-bold">Shared fields</h2>
		<FieldConfigFormShared
			fields={data.normalized_fields}
			onUpdate={(form) => (sharedData = form)}
		/>
	</div>

	<hr />
	{#each Object.entries(data.specific_fields) as [testId, testData]}
		<div class="space-y-4">
			<h2 id={testId} class="text-lg font-bold">
				{testId}
			</h2>
			<FieldConfigForm
				fields={testData.fields}
				jsonConfig={JSON.parse(testData.content)}
				defaultValues={sharedData}
				{defaultFieldsIds}
				onValidUpdate={(form) => {
					$formData[testId] = form;
				}}
			/>
		</div>
		<hr />
	{/each}
</div>

<div class="bg-background/80 sticky bottom-0 flex justify-between border-t p-4 backdrop-blur-lg">
	<div class="flex items-center gap-3">
		{#await completionStatusPromise then [completeTestsCount, incompleteTestsIds]}
			<p>
				{completeTestsCount}/{testsIds.length} tests complete
			</p>
			{#if incompleteTestsIds.length}
				<Popover.Root>
					<Popover.Trigger class="rounded-md p-1 hover:cursor-pointer hover:bg-gray-200">
						{#snippet child({ props })}
							<Button {...props} variant="outline" class="px-3">
								View incomplete tests ({incompleteTestsIds.length})
							</Button>
						{/snippet}
					</Popover.Trigger>
					<Popover.Content class="dark w-fit">
						<ul class="space-y-1 text-sm">
							{#each incompleteTestsIds as testId}
								<li>
									<a class="underline hover:no-underline" href={`#${testId}`}>
										{testId}
									</a>
								</li>
							{/each}
						</ul>
					</Popover.Content>
				</Popover.Root>
			{:else}
				<p>✅</p>
			{/if}
		{/await}
	</div>

	<Form {form} hide={['submit_button']}>
		<SubmitButton>Save</SubmitButton>
	</Form>
</div>
