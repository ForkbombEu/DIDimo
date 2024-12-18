<script lang="ts" generics="C extends CollectionName, Expand extends ExpandQueryOption<C> = never">
	import { type ExpandQueryOption, type QueryResponse, PocketbaseQuery } from '@/pocketbase/query';
	import type { CollectionName } from '@/pocketbase/collections-models';
	import { createRecordDisplay } from './utils';
	import Search from '@/components/ui-custom/search.svelte';
	import type { SearchFunction } from '@/components/ui-custom/search.svelte';
	import type { CollectionInputProps } from './types';

	//

	type Props = CollectionInputProps<C, Expand>;

	let {
		collection,
		queryOptions = {},
		disabled = false,
		label = undefined,
		placeholder = undefined,
		onSelect = () => {},
		displayFields = undefined,
		displayFn = undefined,
		...rest
	}: Props = $props();

	//

	type SearchFn = SearchFunction<QueryResponse<C, Expand>>;

	const searchFunction: SearchFn = $derived(async function (text: string | undefined) {
		const query = new PocketbaseQuery(collection, { ...queryOptions, search: text });
		const records = await query.getFullList();

		return records.map((item) => ({
			value: item,
			label: createRecordDisplay(item, displayFields, displayFn),
			disabled: false
		}));
	});
</script>

<Search searchFn={searchFunction} {onSelect} {label} {placeholder} {disabled} {...rest} />
