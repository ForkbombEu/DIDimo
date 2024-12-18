<script lang="ts" generics="T">
	import { DEFAULT_SORT_ORDER } from '@/pocketbase/query';
	import Icon from '@/components/ui-custom/icon.svelte';
	import { Button } from '@/components/ui/button';
	import { getCollectionManagerContext } from '../collectionManagerContext';
	import { Head } from '@/components/ui/table';
	import type { KeyOf } from '@/utils/types';
	import { capitalize } from 'lodash';
	import { ArrowUp, ArrowDown } from 'lucide-svelte';

	interface Props {
		field: KeyOf<T>;
		label?: string | undefined;
	}

	let { field, label = undefined }: Props = $props();

	const { manager } = $derived(getCollectionManagerContext());

	const sortState = $derived(manager.query.sortOption);
	const isSortField = $derived(sortState[0] == field);

	async function handleClick() {
		if (!isSortField) {
			manager.query.sortBy([field, DEFAULT_SORT_ORDER]);
		} else {
			manager.query.flipSort();
		}
	}
</script>

<Head class="group">
	<div class="flex items-center gap-x-2">
		{label ?? capitalize(field)}
		<Button
			size="icon"
			variant="ghost"
			class="{isSortField ? 'visible' : 'invisible'} size-6 group-hover:visible"
			onclick={handleClick}
		>
			<Icon src={!isSortField ? ArrowUp : sortState[1] == 'DESC' ? ArrowDown : ArrowUp} size={14} />
		</Button>
	</div>
</Head>
