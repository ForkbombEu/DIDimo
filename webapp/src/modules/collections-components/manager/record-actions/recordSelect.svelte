<script lang="ts" generics="C extends CollectionName">
	import type { CollectionResponses } from '@/pocketbase/types';
	import { Checkbox } from '@/components/ui/checkbox';
	import { getCollectionManagerContext } from '../collectionManagerContext';
	import type { CollectionName } from '@/pocketbase/collections-models';

	//

	interface Props {
		record: CollectionResponses[C];
	}

	const { record }: Props = $props();

	//

	const { manager } = $derived(getCollectionManagerContext());
	const checked = $derived(manager.selectedRecords.includes(record.id));

	function handleChecked(checked: boolean | 'indeterminate') {
		if (checked == 'indeterminate') return;
		if (checked) manager.selectRecord(record.id);
		else manager.deselectRecord(record.id);
	}
</script>

<Checkbox {checked} onCheckedChange={handleChecked} value={record.id} name="select" />
