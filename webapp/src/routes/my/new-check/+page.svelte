<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	// import { CollectionForm } from '@/collections-components';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
	import { PageContent } from '@/components/layout';
	import { Badge } from '@/components/ui/badge';

	import { m } from '@/i18n';
</script>

<PageContent>
	<CollectionManager collection="custom_checks" queryOptions={{ expand: ['standard'] }}>
		{#snippet top({ Header })}
			<Header title={m.Custom_checks()} />
		{/snippet}

		{#snippet records({ records, Card })}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each records as record}
					<Card {record} hide={['select', 'share']}>
						{#snippet children({ Title, Description })}
							{@const tags = (record.tags as string).split('#').filter(Boolean)}
							{@const standard = record.expand?.standard}
							<div class="flex flex-col gap-2">
								<Title>{record.name}</Title>
								{#if standard}
									<Badge variant="default" class="w-fit"
										>{standard.name} {standard.version}</Badge
									>
								{/if}
								<Description>{record.description}</Description>
								<div class="flex flex-wrap gap-2">
									{#if record.version}
										<Badge variant="outline">
											{record.version}
										</Badge>
									{/if}
									{#each tags as tag}
										<Badge variant="secondary">{tag}</Badge>
									{/each}
								</div>
							</div>
						{/snippet}
					</Card>
				{/each}
			</div>
		{/snippet}
	</CollectionManager>
</PageContent>
