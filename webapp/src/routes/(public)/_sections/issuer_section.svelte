<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->
<script lang="ts">
	import { m } from '@/i18n';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui-custom/button.svelte';
	import { CollectionManager } from '@/collections-components';
	import RecordCard from '@/collections-components/manager/recordCard.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';

	const MAX_ITEMS = 3;
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<T tag="h3">{m.Find_issuers()}</T>
		<Button variant="default" href="/services">{m.All_issuers()}</Button>
	</div>
	<CollectionManager
		collection="credential_issuers"
		queryOptions={{ perPage: MAX_ITEMS }}
		hide={['pagination']}
	>
		{#snippet records({ records })}
			<PageGrid>
				{#each records as issuer, i}
					{@const isLast = i == MAX_ITEMS - 1}
					<ServiceCard service={issuer} class={isLast ? 'hidden lg:flex' : ''} />
				{/each}
			</PageGrid>
		{/snippet}
	</CollectionManager>
</div>
