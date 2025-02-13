<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<div class="bg-secondary grow">
	<CollectionManager collection="credentials" queryOptions={{ searchFields: ['name'] }}>
		{#snippet top({ Search })}
			<PageTop>
				<T tag="h1">{m.Find_credential_attributes()}</T>

				<Search />
			</PageTop>
		{/snippet}

		{#snippet records({ records })}
			<PageContent class="bg-secondary grow">
				<PageGrid>
					{#each records as record (record.id)}
						<CredentialCard credential={record} />
					{/each}
				</PageGrid>
			</PageContent>
		{/snippet}
	</CollectionManager>
</div>
