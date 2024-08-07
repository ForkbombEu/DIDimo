<script lang="ts">
	import type { RecordFullListOptions } from 'pocketbase';

	import type { RecordInputOptions } from './types';
	import type { Collections } from '$lib/pocketbase/types';
	import { onMount } from 'svelte';
	import type { PBResponse } from '$lib/utils/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { pb } from '$lib/pocketbase';
	import { Select } from 'flowbite-svelte';
	import { createRecordLabel, excludeIdsFilter, mergeFilters } from './utils';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	//

	export let collection: string | Collections;
	export let recordId: string | undefined = undefined;
	export let options: Partial<RecordInputOptions<RecordGeneric>> = {};

	let {
		displayFields = [],
		name = undefined,
		required = false,
		placeholder = undefined,
		filter = undefined,
		expand = undefined,
		formatRecord = undefined
	} = options;

	$: exclude = options.excludeIds ?? [];
	$: disabled = options.disabled ?? false;

	//

	let records: RecordGeneric[] = [];
	$: loadRecords(exclude, filter);

	async function loadRecords(excludeIds: string[], filter: string | undefined) {
		const pbFilter = mergeFilters(excludeIdsFilter(excludeIds), filter);

		const options: RecordFullListOptions = {
			requestKey: null
		};
		if (expand) options.expand = expand;
		if (pbFilter) options.filter = pbFilter;

		records = await pb.collection(collection).getFullList<RecordGeneric>(options);
	}

	onMount(() => {
		pb.collection(collection).subscribe('*', async function (e) {
			loadRecords(exclude, filter);
		});
		return () => {
			pb.collection(collection).unsubscribe('*');
		};
	});

	function createItems(records: RecordGeneric[]) {
		return records.map((r) => ({
			name: formatRecord ? formatRecord(r) : createRecordLabel(r, displayFields),
			value: r.id
		}));
	}

	function handleInput(e: Event) {
		// @ts-ignore
		recordId = e.target.value;
	}
</script>

<Select
	{required}
	{placeholder}
	{disabled}
	{name}
	items={createItems(records)}
	value={recordId}
	on:input={handleInput}
/>
