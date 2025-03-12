<script lang="ts">
	import { WelcomeSession } from '@/auth/welcome';
	import { CollectionManager } from '@/collections-components/index.js';
	import { RecordCard } from '@/collections-components/manager';
	import T from '@/components/ui-custom/t.svelte';
	import { Separator } from '@/components/ui/separator';
	import { currentUser } from '@/pocketbase';

	if (WelcomeSession.isActive()) WelcomeSession.end();
</script>

<div class="flex flex-col space-y-8 p-4">
	<Separator />

	<div class="space-y-2">
		<T tag="h1">Dashboard</T>
		<T>Welcome back, {$currentUser?.name}</T>
	</div>

	<Separator />

	<div class="space-y-2">
		<T tag="h2">Provider claims (in review)</T>
		<CollectionManager
			collection="provider_claims"
			editFormFieldsOptions={{ exclude: ['owner', 'provider', 'status'] }}
		>
			{#snippet records({ records })}
				<div class="grid grid-cols-3 gap-4">
					{#each records as record}
						<RecordCard {record} hide={['select', 'share']}>
							{#snippet children({ Title, Description })}
								<Title>{record.name}</Title>
							{/snippet}
						</RecordCard>
					{/each}
				</div>
			{/snippet}
		</CollectionManager>
	</div>
</div>
