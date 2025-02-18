<script lang="ts">
	import IconButton from '@/components/ui-custom/iconButton.svelte';
	import { getCollectionManagerContext } from './collectionManagerContext';
	import { Input } from '@/components/ui/input';
	import { m } from '@/i18n';
	import { Debounced, watch } from 'runed';
	import { Button } from '@/components/ui/button';
	import Icon from '@/components/ui-custom/icon.svelte';
	import { X } from 'lucide-svelte';
	import { String } from 'effect';

	//

	type Props = {
		class?: string;
	};

	let { class: className }: Props = $props();

	const { manager } = getCollectionManagerContext();

	let searchText = $state('');
	const deboucedSearch = new Debounced(() => searchText, 500);

	$effect(() => {
		manager.query.setSearch(deboucedSearch.current);
	});
</script>

<div class="relative flex">
	<Input bind:value={searchText} placeholder={m.Search()} class={className} />
	{#if String.isString(searchText)}
		<Button
			onclick={() => {
				manager.query.clearSearch();
				searchText = '';
			}}
			class="absolute right-1 top-1 size-8"
			variant="ghost"
		>
			<Icon src={X} size="" />
		</Button>
	{/if}
</div>
