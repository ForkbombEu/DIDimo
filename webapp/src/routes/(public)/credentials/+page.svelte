<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<CollectionManager collection="credentials" queryOptions={{ searchFields: ['name'] }}>
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">{m.Find_credential_attributes()}</T>

			<Search />
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
