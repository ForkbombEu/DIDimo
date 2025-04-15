<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import WalletCard from '$lib/layout/walletCard.svelte';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
	import T from '@/components/ui-custom/t.svelte';
</script>

<CollectionManager collection="wallets">
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">Wallets and verifiers</T>
			<Search class="border-primary bg-secondary" />
		</PageTop>
	{/snippet}

	{#snippet contentWrapper(children)}
		<PageContent class="grow bg-secondary">
			{@render children()}
		</PageContent>
	{/snippet}

	{#snippet records({ records })}
		{console.log(records)}
		<PageGrid>
			{#each records as app}
				<WalletCard {app} />
			{/each}
		</PageGrid>
	{/snippet}
</CollectionManager>
