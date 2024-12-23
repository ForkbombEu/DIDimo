<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import { CollectionManager } from '@/collections-components';
	import CollectionManagerHeader from '@/collections-components/manager/collectionManagerHeader.svelte';
	import { m } from '@/i18n';
	import type { PageData } from './$types';
	import { RecordCard } from '@/collections-components/manager';

	let { data }: { data: PageData } = $props();
</script>

<PageContent>
	<CollectionManager
		collection="services"
		formFieldsOptions={{
			hide: {
				organization: data.organization.id
			}
		}}
	>
		{#snippet top()}
			<CollectionManagerHeader title={m.Services()} />
		{/snippet}

		{#snippet records({ records })}
			<div class="mt-4">
				{#each records as record}
					<RecordCard {record} hide={['select', 'share']}>
						{#snippet children({ Title, Description })}
							<Title>{record.name}</Title>
							<Description>{record.description}</Description>
						{/snippet}
					</RecordCard>
				{/each}
			</div>
		{/snippet}
	</CollectionManager>
</PageContent>
