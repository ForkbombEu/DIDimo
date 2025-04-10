<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<CollectionManager
	collection="credential_issuers"
	queryOptions={{ expand: ['credentials_via_credential_issuer'] }}
>
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">{m.Find_identity_solutions()}</T>
			<Search class="border-primary bg-secondary" />
		</PageTop>
	{/snippet}

	{#snippet contentWrapper(children)}
		<PageContent class="bg-secondary grow">
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
