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

	const MAX_ITEMS = 3;
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<T tag="h3">{m.Find_verifiers()}</T>
		<Button variant="default" href="/services">{m.All_verifiers()}</Button>
	</div>
	<CollectionManager
		collection="verifiers"
		queryOptions={{ perPage: MAX_ITEMS }}
		hide={['pagination']}
	>
		{#snippet records({ records })}
			<PageGrid>
				{#each records as credential, i}
					{@const isLast = i == MAX_ITEMS - 1}
					<RecordCard record={credential} class={isLast ? 'hidden lg:flex' : ''} />
				{/each}
			</PageGrid>
		{/snippet}
	</CollectionManager>
</div>
