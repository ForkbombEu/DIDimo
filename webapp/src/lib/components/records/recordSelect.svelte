<script lang="ts">
	import type { RecordInputOptions } from './types';
	import type { Collections } from '$lib/pocketbase/types';
	import { onMount } from 'svelte';
	import type { PBResponse } from '$lib/utils/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { pb } from '$lib/pocketbase';
	import { Select } from 'flowbite-svelte';
	import { createRecordLabel, excludeStringArray } from './utils';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	//

	export let collection: string | Collections;
	export let recordId: string | undefined = undefined;
	export let options: Partial<RecordInputOptions<RecordGeneric>> = {};

	let { displayFields = [], name = undefined, required = false, placeholder = undefined } = options;
	$: exclude = options.excludeIds ?? [];
	$: disabled = options.disabled ?? false;

	//

	let records: RecordGeneric[] = [];
	$: loadRecords(exclude);

	async function loadRecords(excludeIds: string[]) {
		records = await pb.collection(collection).getFullList<RecordGeneric>({
			requestKey: null,
			filter: excludeStringArray('id', excludeIds)
		});
	}

	onMount(() => {
		pb.collection(collection).subscribe('*', async function (e) {
			loadRecords(exclude);
		});
		return () => {
			pb.collection(collection).unsubscribe('*');
		};
	});

	function createItems(records: RecordGeneric[]) {
		return records.map((r) => ({
			name: createRecordLabel(r, displayFields),
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
