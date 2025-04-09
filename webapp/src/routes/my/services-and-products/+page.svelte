<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { CollectionManager } from '@/collections-components';
	import { m } from '@/i18n';
	import { Plus } from 'lucide-svelte';
	import * as Dialog from '@/components/ui/dialog';
	import { buttonVariants } from '@/components/ui/button';
	import CredentialIssuerForm from './credential-issuer-form.svelte';
	import T from '@/components/ui-custom/t.svelte';

	//

	let isCredentialIssuerModalOpen = $state(false);
</script>

<div class="space-y-12">
	<div class="space-y-4">
		<CollectionManager collection="credential_issuers">
			{#snippet top({ Header })}
				<Header title={m.Credential_issuers()}>
					{#snippet right()}
						{@render CreateCredentialIssuerModal()}
					{/snippet}
				</Header>
			{/snippet}

			{#snippet records({ records, Card })}
				{#each records as record}
					<Card {record} hide={['select', 'share']} class="bg-background">
						<T>{record.url}</T>
					</Card>
				{/each}
			{/snippet}
		</CollectionManager>
	</div>

	<div class="space-y-4">
		<CollectionManager collection="wallets">
			{#snippet top({ Header })}
				<Header title={m.Wallets()} />
			{/snippet}
		</CollectionManager>
	</div>
</div>

<!--  -->

{#snippet CreateCredentialIssuerModal()}
	<Dialog.Root bind:open={isCredentialIssuerModalOpen}>
		<Dialog.Trigger class={buttonVariants({ variant: 'default' })}>
			<Plus />
			{m.Add_new_credential_issuer()}
		</Dialog.Trigger>

		<Dialog.Content class=" sm:max-w-[425px]">
			<Dialog.Header>
				<Dialog.Title>{m.Add_new_credential_issuer()}</Dialog.Title>
			</Dialog.Header>

			<div class="pt-8">
				<CredentialIssuerForm
					onSuccess={() => {
						isCredentialIssuerModalOpen = false;
					}}
				/>
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/snippet}
