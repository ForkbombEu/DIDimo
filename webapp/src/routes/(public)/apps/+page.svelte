<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import CardLink from '$lib/layout/card-link.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import Card from '@/components/ui-custom/card.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
</script>

<CollectionManager collection="wallets">
	{#snippet top({ Search })}
		<PageTop>
			<T tag="h1">Wallets and verifiers</T>
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
			{#each records as organization}
				<CardLink href={`/organizations/${organization.id}`}>
					<div class="flex items-center gap-2">
						{#if organization.logo}
							<Avatar
								src={organization.logo}
								class="!rounded-sm"
								hideIfLoadingError
							/>
						{/if}
						<T class="font-semibold">{organization.name}</T>
					</div>
				</CardLink>
			{/each}
		</PageGrid>
	{/snippet}
</CollectionManager>
