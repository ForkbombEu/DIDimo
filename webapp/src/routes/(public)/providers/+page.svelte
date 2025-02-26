<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import * as Sheet from '@/components/ui/sheet';
</script>

<CollectionManager collection="services" queryOptions={{ expand: ['credential_issuers'] }}>
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">{m.Find_identity_solutions()}</T>
			<Search class="border-primary bg-secondary" />

			<Sheet.Root>
				<Sheet.Trigger>Open</Sheet.Trigger>
				<Sheet.Content>
					<Sheet.Header>
						<Sheet.Title>Are you sure absolutely sure?</Sheet.Title>
						<Sheet.Description>
							This action cannot be undone. This will permanently delete your account
							and remove your data from our servers.
						</Sheet.Description>
					</Sheet.Header>
				</Sheet.Content>
			</Sheet.Root>
		</PageTop>
	{/snippet}

	{#snippet contentWrapper(children)}
		<PageContent class="grow bg-secondary">
			{@render children()}
		</PageContent>
	{/snippet}

	{#snippet records({ records })}
		<PageGrid>
			{#each records as service}
				<ServiceCard {service} class="grow" />
			{/each}
		</PageGrid>
	{/snippet}
</CollectionManager>
