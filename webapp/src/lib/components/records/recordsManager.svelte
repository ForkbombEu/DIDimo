<script lang="ts" context="module">
	import type { RecordInputOptions } from './types';
	import type { PBResponse } from '$lib/utils/types';

	export type InputMode = 'search' | 'select';

	export type RecordManagerActions = 'edit' | 'create';

	export type RecordsManagerOptions<R extends PBResponse = PBResponse> = Partial<{
		inputMode: InputMode;
		multiple: boolean;
		showActions: RecordManagerActions[];
		max: number | undefined;
		formSettings: Partial<FieldsSettings<R>>;
	}> &
		Partial<RecordInputOptions<R>>;
</script>

<script lang="ts">
	import { pb } from '$lib/pocketbase';

	import RecordSelect from './recordSelect.svelte';
	import RecordSearch from './recordSearch.svelte';
	import ArrayOrItemManager from '$lib/components/arrayOrItemManager.svelte';

	import { createRecordLabel } from './utils';
	import { createTypeProp } from '$lib/utils/typeProp';

	import Drawer from '$lib/components/drawer.svelte';
	import { Button } from 'flowbite-svelte';
	import { Plus } from 'svelte-heros-v2';
	import { createToggleStore } from '../utils/toggleStore';
	import RecordForm, { type FieldsSettings } from '$lib/recordForm/recordForm.svelte';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	export let collection: string;
	export let value: string[] | string | undefined = undefined;
	export let options: Partial<RecordsManagerOptions<RecordGeneric>> = {};

	let {
		inputMode = 'search',
		multiple = false,
		displayFields = [],
		excludeIds = [],
		showActions = [],
		max = undefined,
		formSettings = {},
		formatRecord = undefined
	} = options;

	//

	let inputComponent: typeof RecordSearch<RecordGeneric> | typeof RecordSelect<RecordGeneric>;
	$: if (inputMode == 'search') inputComponent = RecordSearch;
	else if (inputMode == 'select') inputComponent = RecordSelect;

	//

	let tempId: string | undefined = undefined;
	$: handleSelect(tempId);

	function handleSelect(id: typeof tempId) {
		if (!id) return;

		if (!multiple) {
			value = id;
		} else {
			if (!value) value = [id];
			else {
				if (!value.includes(id)) {
					value = [...value, id];
				}
			}
		}

		tempId = undefined;
	}

	//

	let tempIds: string[] = [];
	$: setupTempIds(value);

	function setupTempIds(v: typeof value) {
		if (Array.isArray(v)) tempIds = v;
		else if (v) tempIds = [v];
		else tempIds = [];
	}

	//

	$: exclude = [...excludeIds, ...tempIds];
	$: disabled = isDisabled(value);

	function isDisabled(v: typeof value): boolean {
		if (!v || !Array.isArray(v) || !max) return false;
		return v.length >= max;
	}

	//

	let tempRecords: Record<string, RecordGeneric> = {};
	$: loadDisplayRecords(tempIds);

	async function loadDisplayRecords(ids: typeof tempIds) {
		const records = await pb.collection(collection).getFullList<RecordGeneric>({
			filter: ids.map((id) => `id = '${id}'`).join(' || '),
			requestKey: null
		});
		tempRecords = {};
		records.forEach((r) => {
			tempRecords[r.id] = r;
		});
	}

	//

	const hideCreateDrawer = createToggleStore(true);
</script>

<div class="space-y-4">
	<svelte:component
		this={inputComponent}
		{recordType}
		{collection}
		bind:recordId={tempId}
		options={{ ...options, excludeIds: exclude, disabled }}
	/>

	<ArrayOrItemManager bind:value let:item>
		{@const record = tempRecords[item]}
		{#if record}
			<slot {record}>
				{formatRecord ? formatRecord(record) : createRecordLabel(record, displayFields)}
			</slot>
		{/if}
	</ArrayOrItemManager>

	{#if showActions.includes('create')}
		<div class="flex justify-end">
			<Button color="alternative" size="xs" on:click={hideCreateDrawer.off}>
				<Plus size="16" />
				<span class="ml-1">Create new record</span>
			</Button>
		</div>
	{/if}
</div>

<Drawer bind:hidden={$hideCreateDrawer} placement="right" width="w-3/6" title="Create new record">
	<div class="p-6">
		<RecordForm
			{collection}
			fieldsSettings={formSettings}
			on:create={(e) => {
				tempId = e.detail.record.id;
				hideCreateDrawer.on();
			}}
		/>
	</div>
</Drawer>
