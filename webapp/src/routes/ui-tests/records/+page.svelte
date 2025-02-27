<script lang="ts">
	import { CollectionSelect } from '@/collections-components';
	import CollectionSearch from '@/collections-components/collectionSearch.svelte';
</script>

<div class="space-y-4">
	<CollectionSearch
		collection="z_test_collection"
		onSelect={(record) => console.log(record)}
		queryOptions={{ expand: ['relation_multi_field'], searchFields: ['text_field'] }}
		displayFn={(record) => record.text_field}
	/>

	<CollectionSelect
		collection="z_test_collection"
		onSelect={(r) => console.log(r)}
		queryOptions={{ expand: ['relation_multi_field'] }}
		displayFn={(record) => {
			return record.expand?.relation_multi_field?.map((u) => u.name).join(', ') ?? '';
		}}
		clearValueOnSelect
	/>
</div>
