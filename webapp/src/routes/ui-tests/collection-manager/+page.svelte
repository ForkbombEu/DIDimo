<script lang="ts">
	import { CollectionManager, type FilterGroup } from '@/collections-components/manager';

	const filters: FilterGroup[] = [
		{
			name: 'Group 1',
			filters: [
				{
					name: 'Has self relation',
					id: 'has_self_relation',
					expression: 'self_relation != null'
				},
				{
					name: 'Self relation has Select B',
					id: 'has_relation_of_select_b',
					expression: "self_relation.select_field = 'b'"
				}
			]
		}
	];
</script>

<CollectionManager
	collection="z_test_collection"
	{filters}
	queryOptions={{
		expand: ['relation_field'],
		perPage: 6,
		sort: ['self_relation', 'ASC']
	}}
>
	{#snippet top({ Header, Search, Filters })}
		<Header />

		<div class="mb-4 mt-4">
			<Search />
			<Filters />
		</div>
	{/snippet}

	{#snippet records({ records, Table, Card })}
		<Table {records} fields={['id', 'text_field', 'self_relation', 'select_field']}></Table>

		<div class="mt-4 space-y-2">
			{#each records as record}
				<Card {record}>
					{#snippet children({ Title, Description })}
						<Title>{record.text_field}</Title>
						<Description>{record.expand?.relation_field?.email}</Description>
					{/snippet}
				</Card>
			{/each}
		</div>
	{/snippet}
</CollectionManager>
