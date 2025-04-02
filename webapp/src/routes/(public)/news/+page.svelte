<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import NewsCard from '$lib/layout/newsCard.svelte';
</script>

<CollectionManager
	collection="news"
	queryOptions={{
		searchFields: ['title', 'summary', 'news'],
		perPage: 20
	}}
>
	{#snippet top({ Search, Filters })}
		<PageTop>
			<T tag="h1">{m.News()}</T>
			<div class="flex items-center gap-2">
				<Search class="border-primary bg-secondary" containerClass="grow" />
				<Filters>
					{m.filters()}
				</Filters>
			</div>
		</PageTop>
	{/snippet}

	{#snippet contentWrapper(children)}
		<PageContent class="grow bg-secondary">
			{@render children()}
		</PageContent>
	{/snippet}

	{#snippet records({ records })}
		<PageGrid>
			{#each records as record (record.id)}
				<NewsCard news={record} />
			{/each}
		</PageGrid>
	{/snippet}
</CollectionManager>
