<script lang="ts" generics="C extends CollectionName, Expand extends ExpandQueryOption<C> = never">
	import { type ExpandQueryOption, PocketbaseQuery, type QueryResponse } from '@/pocketbase/query';
	import type { CollectionName } from '@/pocketbase/collections-models';
	import { setupComponentPocketbaseSubscriptions } from '@/pocketbase/subscriptions';
	import type { RecordIdString } from '@/pocketbase/types';
	import { createRecordDisplay } from './utils';
	import SelectInput, { type SelectItem } from '@/components/ui-custom/selectInput.svelte';
	import type { CollectionInputProps } from './types';

	//

	type Props = CollectionInputProps<C, Expand>;

	let {
		collection,
		queryOptions = {},
		disabled = false,
		placeholder,
		clearValueOnSelect = false,
		onSelect = () => {},
		displayFields,
		displayFn,
		controlAttrs
	}: Props = $props();

	//

	type Record = QueryResponse<C, Expand>;

	let records = $state<Record[]>([]);
	let recordId = $state<RecordIdString | undefined>();
	const selectedRecord = $derived(records.find((r) => r.id == recordId));

	const presentRecord = $derived(function (record: Record) {
		return createRecordDisplay(record, displayFields, displayFn);
	});

	const selectItems: SelectItem[] = $derived(
		records.map((r) => ({
			value: r.id,
			label: presentRecord(r)
		}))
	);

	//

	const loadRecords = $derived(function () {
		const query = new PocketbaseQuery(collection, queryOptions);
		query.getFullList().then((res) => (records = res));
	});

	$effect(() => {
		loadRecords();
	});

	setupComponentPocketbaseSubscriptions({
		collection,
		callback: () => loadRecords(),
		expandOption: queryOptions.expand,
		other: ['authorizations']
	});

	//

	$effect(() => {
		if (!selectedRecord) return;
		onSelect($state.snapshot(selectedRecord) as Record);
		if (clearValueOnSelect) recordId = undefined;
	});
</script>

<SelectInput
	type="single"
	items={selectItems}
	bind:value={recordId}
	{placeholder}
	{controlAttrs}
	{disabled}
>
	{#snippet trigger()}
		{#if selectedRecord}
			{presentRecord(selectedRecord)}
		{/if}
	{/snippet}
</SelectInput>
