<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { CollectionForm } from '@/collections-components/index.js';
	import { PageCard } from '@/components/layout/index.js';
	import Alert from '@/components/ui-custom/alert.svelte';
	import Button from '@/components/ui-custom/button.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { currentUser } from '@/pocketbase/index.js';
	import { InfoIcon, Plus, Pencil } from 'lucide-svelte';
	import OrganizationPageDemo from '$lib/pages/organization-page.svelte';

	let { data } = $props();
	const { organizationInfo } = $derived(data);

	//

	let showOrganizationForm = $state(false);
</script>

{#if !showOrganizationForm}
	{#if organizationInfo}
		<div class="mb-6 flex items-center justify-between">
			<T tag="h3">Page preview</T>
			<Button onclick={() => (showOrganizationForm = true)}>
				<Pencil />
				Edit organization info
			</Button>
		</div>
		<div class="rounded-lg bg-white p-6">
			<OrganizationPageDemo {organizationInfo} />
		</div>
	{:else}
		<Alert variant="info" icon={InfoIcon}>
			{#snippet content({ Title, Description })}
				<Title>Info</Title>
				<Description class="mt-2">
					An organization page is used to present your services and products on the
					marketplace. Create one to get started!
				</Description>
				<div class="mt-2 flex justify-end">
					<Button
						variant="outline"
						onclick={() => {
							showOrganizationForm = true;
						}}
					>
						<Plus />
						Create organization page
					</Button>
				</div>
			{/snippet}
		</Alert>
	{/if}
{:else}
	<PageCard>
		<T tag="h3">
			{organizationInfo ? 'Update your organization page' : 'Create your organization page'}
		</T>
		<!-- TODO - Set owner via hook -->
		<CollectionForm
			collection="organization_info"
			fieldsOptions={{ hide: { owner: $currentUser!.id } }}
			uiOptions={{ hideRequiredIndicator: true }}
			onSuccess={() => {
				invalidateAll();
				showOrganizationForm = false;
			}}
			initialData={organizationInfo}
			recordId={organizationInfo?.id}
		>
			{#snippet submitButtonContent()}
				{organizationInfo ? 'Update organization page' : 'Create organization page'}
			{/snippet}
		</CollectionForm>
	</PageCard>
{/if}
