<script lang="ts">
	import IconButton from '@/components/ui-custom/iconButton.svelte';
	import { getCollectionManagerContext } from './collectionManagerContext';
	import { Input } from '@/components/ui/input';
	import { m } from '@/i18n';
	import { Debounced } from 'runed';

	const { manager } = $derived(getCollectionManagerContext());

	let searchText = $state('');
	const deboucedSearch = new Debounced(() => searchText, 500);
	$effect(() => manager.search(deboucedSearch.current));
</script>

<div class="flex gap-2">
	<Input bind:value={searchText} placeholder={m.Search()} />
	<IconButton
		onclick={() => {
			manager.clearSearch();
			searchText = '';
		}}
	/>
</div>
