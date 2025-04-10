<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { CollectionForm, CollectionManager } from '@/collections-components';
	import { m } from '@/i18n';
	import { InfoIcon, Pencil, Plus } from 'lucide-svelte';
	import * as Dialog from '@/components/ui/dialog';
	import { buttonVariants } from '@/components/ui/button';
	import CredentialIssuerForm from './credential-issuer-form.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { String } from 'effect';
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui-custom/button.svelte';
	import EditCredentialDialog from './edit-credential-dialog.svelte';
	import Alert from '@/components/ui-custom/alert.svelte';

	//

	let { data } = $props();
	const { organizationInfo } = $derived(data);

	let isCredentialIssuerModalOpen = $state(false);
</script>

<div class="space-y-12">
	{#if !organizationInfo}
		<Alert variant="warning" icon={InfoIcon}>
			{#snippet content({ Title, Description })}
				<Title>Important!</Title>
				<Description class="mt-2">
					Before effectively publishing your services and products to the marketplace, you
					need to create a public organization page.
				</Description>
				<div class="mt-2 flex justify-end">
					<Button
						variant="outline"
						href="/my/organization-page"
						onclick={() => {
							isCredentialIssuerModalOpen = true;
						}}
					>
						<Plus />
						Create organization page
					</Button>
				</div>
			{/snippet}
		</Alert>
	{/if}

	<div class="space-y-4">
		<CollectionManager
			collection="credential_issuers"
			queryOptions={{ expand: ['credentials_via_credential_issuer'] }}
			editFormFieldsOptions={{ exclude: ['owner', 'url'], order: ['published'] }}
			subscribe="expanded_collections"
		>
			{#snippet top({ Header })}
				<Header title={m.Credential_issuers()}>
					{#snippet right()}
						{@render CreateCredentialIssuerModal()}
					{/snippet}
				</Header>
			{/snippet}

			{#snippet records({ records, Card })}
				{#each records as record}
					{@const credentials = record.expand?.credentials_via_credential_issuer ?? []}
					<Card {record} hide={['select', 'share']} class="bg-background">
						{@const title = String.isNonEmpty(record.name) ? record.name : record.url}
						<T class="font-bold">
							{title}
						</T>
						{#if title != record.url}
							<T class="">
								{record.url}
							</T>
						{/if}
						{#if record.published}
							<Badge variant="default">{m.Published()}</Badge>
						{/if}

						{#if credentials.length === 0}
							<T>{m.No_credentials_available()}</T>
						{:else}
							<T>{m.count_available_credentials({ number: credentials.length })}</T>

							<ul>
								{#each credentials as credential}
									<li>
										{credential.key}
										{#if credential.published}
											<Badge variant="default">{m.Published()}</Badge>
										{/if}

										<EditCredentialDialog {credential} />
									</li>
								{/each}
							</ul>
						{/if}
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
