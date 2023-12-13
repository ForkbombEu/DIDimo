<script lang="ts">
	import RecordSelect from '$lib/components/records/recordSelect.svelte';
	// import RelationsManager from '$lib/components/relationsManager.svelte';
	import { Collections } from '$lib/pocketbase/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import type {
		CrudExampleRecord,
		CrudExampleResponse,
		FeaturesResponse,
		UsersResponse
	} from '$lib/pocketbase/types';
	import RecordSearch from '$lib/components/records/recordSearch.svelte';
	import RecordsManager from '$lib/components/records/recordsManager.svelte';
	import { boolean } from 'zod';

	let searchValue: string | undefined = undefined;
	const collection = Collections.CrudExample;
	const recordType =
		createTypeProp<CrudExampleResponse<{ owner: UsersResponse; relation: FeaturesResponse[] }>>();
</script>

<div class="p-8 space-y-4">
	<p>Record select</p>
	<RecordSelect {recordType} {collection} bind:recordId={searchValue} />
	<pre>{JSON.stringify(searchValue)}</pre>
</div>

<div class="p-8 space-y-4">
	<p>Record search</p>
	<RecordSearch {collection} bind:recordId={searchValue} options={{ placeholder: 'mimmo' }} />
</div>

<div class="p-8 space-y-4">
	<p>Records manager</p>

	<RecordsManager
		{collection}
		{recordType}
		options={{
			multiple: true,
			inputMode: 'select',
			max: 3,
			showActions: ['create'],
			displayFields: ['id', 'text'],
			formSettings: {
				order: ['image', 'multiselect'],
				hide: {
					text: 'ao'
				},
				relations: {
					owner: {
						displayFields: ['id', 'email']
					},
					file_only_pdf_json: {
						displayFields: ['sad']
					},
					relation: {
						displayFields: ['name', 'id', 'created']
					}
				}
			}
		}}
	/>
</div>
