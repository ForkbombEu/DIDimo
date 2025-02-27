<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
	import * as Sheet from '@/components/ui/sheet/index.js';
</script>

<CollectionManager collection="credentials" queryOptions={{ searchFields: ['name'] }}>
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">{m.Find_credential_attributes()}</T>

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
		<PageContent class="bg-secondary grow">
			{@render children()}
		</PageContent>
	{/snippet}

	{#snippet records({ records })}
		<PageGrid>
			{#each records as record (record.id)}
				<CredentialCard credential={record} />
			{/each}
		</PageGrid>
	{/snippet}
</CollectionManager>
