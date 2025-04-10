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
	import { RecordDelete, RecordEdit } from '@/collections-components/manager';
	import PublishButton from './publish-button.svelte';
	import Separator from '@/components/ui/separator/separator.svelte';
	import { currentUser } from '@/pocketbase';

	//

	let { data } = $props();
	const { organizationInfo } = $derived(data);
	const canPublish = $derived(Boolean(organizationInfo));

	let isCredentialIssuerModalOpen = $state(false);
</script>

<div class="space-y-12">
	<div class="space-y-4">
		<CollectionManager
			collection="credential_issuers"
			queryOptions={{
				expand: ['credentials_via_credential_issuer'],
				filter: `owner.id = "${$currentUser?.id}"`
			}}
			editFormFieldsOptions={{ exclude: ['owner', 'url', 'published'], order: ['published'] }}
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
					<Card
						{record}
						hide={['select', 'share', 'delete', 'edit']}
						class="bg-background"
					>
						{@const title = String.isNonEmpty(record.name) ? record.name : record.url}
						<div class="space-y-4">
							<div class="flex items-center justify-between gap-4">
								<div>
									<div class="flex items-center gap-2">
										<T class="font-bold">
											{title}
										</T>
										{#if record.published}
											<Badge variant="default">{m.Published()}</Badge>
										{/if}
									</div>

									{#if title != record.url}
										<T class="">
											{record.url}
										</T>
									{/if}
								</div>

								<div class="flex items-center gap-1">
									<PublishButton {record} {canPublish} {cannotPublishMessage}>
										{#snippet button({ togglePublish, label })}
											<Button variant="outline" onclick={togglePublish}>
												{label}
											</Button>
										{/snippet}
									</PublishButton>

									<RecordEdit {record} />
									<RecordDelete {record} />
								</div>
							</div>

							<Separator />

							{#if credentials.length === 0}
								<T class="text-gray-300">{m.No_credentials_available()}</T>
							{:else}
								<T>
									{m.count_available_credentials({
										number: credentials.length
									})}
								</T>

								<ul class="space-y-2">
									{#each credentials as credential}
										<li
											class="bg-muted flex items-center justify-between rounded-md p-2 px-4"
										>
											<div class="flex items-center gap-2">
												{credential.key}
												{#if credential.published}
													<Badge variant="default">{m.Published()}</Badge>
												{/if}
											</div>

											<div class="flex items-center gap-1">
												{#if record.published}
													<PublishButton
														record={credential}
														{cannotPublishMessage}
														canPublish={record.published && canPublish}
													>
														{#snippet button({ togglePublish, label })}
															<Button
																variant="outline"
																onclick={togglePublish}
																class="h-fit py-1 text-xs"
															>
																{label}
															</Button>
														{/snippet}
													</PublishButton>
												{/if}
												<EditCredentialDialog {credential} />
											</div>
										</li>
									{/each}
								</ul>
							{/if}
						</div>
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

<!--  -->

{#snippet cannotPublishMessage()}
	<T>
		Before you can publish your service, you need to create a public profile for your
		organization.
	</T>
	<div class="flex justify-end">
		<Button href="/my/organization-page">
			<Plus />
			{m.Create_organization()}
		</Button>
	</div>
{/snippet}
