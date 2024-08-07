<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { getRecordsManagerContext } from '../../collectionManager.svelte';
	import type { PBResponse } from '$lib/utils/types';

	import { RecordForm, type FieldsSettings } from '$lib/recordForm';
	import { Button, Modal } from 'flowbite-svelte';
	import PortalWrapper from '$lib/components/portalWrapper.svelte';
	import { Plus } from 'svelte-heros-v2';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	export let initialData: Partial<RecordGeneric> = {};
	export let formSettings: Partial<FieldsSettings<RecordGeneric>> = {};

	//

	const dispatch = createEventDispatcher<{
		success: {
			record: RecordGeneric;
		};
	}>();

	let { collection, dataManager, formFieldsSettings } = getRecordsManagerContext<RecordGeneric>();
	let { base, create } = formFieldsSettings;
	let fieldsSettings = { ...base, ...create, ...formSettings };
	let { loadRecords } = dataManager;

	//

	let open = false;

	function openModal() {
		open = true;
	}
</script>

<slot {openModal}>
	<Button color="alternative" on:click={openModal}>
		<Plus size="20" />
		<span class="ml-1">Add entry</span>
	</Button>
</slot>

<PortalWrapper>
	<Modal bind:open title="Create record" size="md">
		<div class="w-full">
			<RecordForm
				{collection}
				{fieldsSettings}
				{initialData}
				on:success={async (e) => {
					await loadRecords();
					dispatch('success', { record: e.detail.record });
					open = false;
				}}
			/>
		</div>
	</Modal>
</PortalWrapper>
