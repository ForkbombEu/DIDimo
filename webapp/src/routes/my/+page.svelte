<script lang="ts">
	import { WelcomeSession } from '@/auth/welcome';
	import { CollectionManager } from '@/collections-components/index.js';
	import { RecordCard } from '@/collections-components/manager';
	import T from '@/components/ui-custom/t.svelte';
	import { Separator } from '@/components/ui/separator';
	import { currentUser } from '@/pocketbase';

	if (WelcomeSession.isActive()) WelcomeSession.end();

	let { data } = $props();
</script>

<div class="flex flex-col space-y-8 p-4">
	<Separator />

	<div class="space-y-2">
		<T tag="h1">Dashboard</T>
		<T>Welcome back, {$currentUser?.name}</T>
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
