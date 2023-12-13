<script lang="ts">
	import { Collections, type CrudExampleResponse } from '$lib/pocketbase/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import {
		CollectionManager,
		CollectionSearch,
		CollectionManagerHeader,
		CollectionEmptyState,
		CollectionTable,
		RecordCard
	} from '$lib/collectionManager';
	import Chip from '$lib/collectionManager/ui/fieldComponents/chip.svelte';

	import { Heading, Hr } from 'flowbite-svelte';
	import Pagination from '$lib/collectionManager/ui/pagination.svelte';
	import EditRecord from '$lib/collectionManager/ui/recordActions/editRecord.svelte';

	const recordType = createTypeProp<CrudExampleResponse>();
</script>

<div class="p-4">
	<CollectionManager
		{recordType}
		collection={Collections.CrudExample}
		formSettings={{
			labels: {
				text: 'Titolo',
				file_only_pdf_json: 'PDF or JSON file',
				textarea: 'Short description'
			},
			relations: {
				relation: {
					displayFields: ['name']
				}
			}
		}}
		editFormSettings={{
			exclude: ['select']
		}}
		let:records
	>
		<div class="space-y-8">
			<CollectionManagerHeader />

			<Hr />

			<div class="space-y-4">
				<Heading tag="h4">Table</Heading>

				<CollectionSearch {recordType} searchableFields={['text', 'textarea']} />
				<CollectionTable
					{records}
					fields={['id', 'text', 'textarea']}
					fieldsLabels={{ id: 'Unique ID', text: 'Title' }}
					hideActions={['edit', 'delete']}
				>
					<svelte:fragment slot="emptyState">
						<CollectionEmptyState
							title="Custom empty state"
							description="Displaying some custom text here!"
						/>
					</svelte:fragment>
					<svelte:fragment slot="actions" let:record>
						<EditRecord {record} formSettings={{ exclude: ['relation_single'] }} />
					</svelte:fragment>
				</CollectionTable>
				<Pagination />
			</div>

			<Hr />

			<div class="space-y-4">
				<Heading tag="h4">Cards</Heading>
				{#if records.length === 0}
					<CollectionEmptyState />
				{:else}
					<div class="grid grid-cols-4 gap-4">
						{#each records as record}
							<RecordCard {record} titleField="id" fieldsComponents={{ select: Chip }} />
						{/each}
					</div>
					<Pagination />
				{/if}
			</div>
		</div>
	</CollectionManager>
</div>
