<script lang="ts">
	import { CollectionManager } from '@/collections-components/manager';
</script>

<CollectionManager
	collection="z_test_collection"
	queryOptions={{
		expand: ['relation_field'],
		perPage: 6
	}}
>
	{#snippet top({ Header, Search })}
		<Header />

		<div class="mb-4 mt-4">
			<Search />
		</div>
	{/snippet}

	{#snippet records({ records, Table, Card })}
		<Table {records} fields={['id', 'text_field', 'relation_field', 'created']}></Table>

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
