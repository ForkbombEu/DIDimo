<script lang="ts">
	import SectionTitle from '@/components/ui-custom/sectionTitle.svelte';
	import RecordCreate from './record-actions/recordCreate.svelte';
	import { getCollectionManagerContext } from './collectionManagerContext';

	interface Props {
		title?: string | undefined;
		hideCreate?: boolean;
	}

	const { title, hideCreate = false }: Props = $props();
	const { manager } = $derived(getCollectionManagerContext());
</script>

<SectionTitle title={title ?? manager.collection}>
	{#snippet right()}
		{#if !hideCreate}
			<RecordCreate />
		{/if}
	{/snippet}
</SectionTitle>
