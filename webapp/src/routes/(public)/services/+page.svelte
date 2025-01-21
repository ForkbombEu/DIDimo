<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<div class="bg-secondary grow">
	<CollectionManager collection="services" queryOptions={{ expand: ['credential_issuers'] }}>
		{#snippet top({ Search })}
			<PageTop>
				<T tag="h1">{m.Find_identity_solutions()}</T>
				<Search />
			</PageTop>
		{/snippet}

		{#snippet records({ records })}
			<PageContent>
				<PageGrid>
					{#each records as service}
						<ServiceCard {service} class="grow" />
					{/each}
				</PageGrid>
			</PageContent>
		{/snippet}
	</CollectionManager>
</div>
