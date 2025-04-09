<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { WelcomeSession } from '@/auth/welcome';
	import { CollectionManager } from '@/collections-components/index.js';
	import { RecordCard } from '@/collections-components/manager';
	import Button from '@/components/ui-custom/button.svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Separator } from '@/components/ui/separator';
	import { m } from '@/i18n/index.js';
	import { currentUser } from '@/pocketbase';
	import { Sparkle, Workflow } from 'lucide-svelte';

	if (WelcomeSession.isActive()) WelcomeSession.end();

	let { data } = $props();
</script>

<div class="flex flex-col space-y-8 p-4">
	<Separator />

	<div class="flex items-center justify-between">
		<div class="space-y-2">
			<T tag="h1">Dashboard</T>
			<T>Welcome back, {$currentUser?.name}</T>
		</div>

		<div class="flex items-center gap-2">
			<Button href="/my/tests/runs" variant="outline">
				<Icon src={Workflow} />
				{m.View_test_runs()}
			</Button>
			<Button href="/my/tests">
				<Icon src={Sparkle} />
				Start a new check
			</Button>
		</div>
	</div>

	<Separator />

	<div class="grid gap-8 pb-8 md:grid-cols-2">
		<div class="space-y-2">
			<T tag="h2">Your providers</T>
			<CollectionManager
				collection="services"
				emptyStateDescription="Start a test to add a provider"
				queryOptions={{
					filter: `owner = "${data.organization?.id}"`
				}}
				editFormFieldsOptions={{ exclude: ['owner', 'wallets', 'credential_issuers'] }}
			>
				{#snippet records({ records })}
					{#each records as record}
						<RecordCard {record} hide={['select', 'share', 'delete']}>
							{#snippet children({ Title, Description })}
								<Title>{record.name}</Title>
							{/snippet}
						</RecordCard>
					{/each}
				{/snippet}
			</CollectionManager>
		</div>

		<div class="space-y-2">
			<T tag="h2">Provider claims (in review)</T>
			<CollectionManager
				collection="provider_claims"
				editFormFieldsOptions={{ exclude: ['owner', 'provider', 'status'] }}
				emptyStateDescription="Claims in review will appear here"
			>
				{#snippet records({ records })}
					{#each records as record}
						<RecordCard {record} hide={['select', 'share']}>
							{#snippet children({ Title, Description })}
								<Title>{record.name}</Title>
							{/snippet}
						</RecordCard>
					{/each}
				{/snippet}
			</CollectionManager>
		</div>
	</div>
</div>
