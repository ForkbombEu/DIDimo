<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageHeader from '$lib/layout/pageHeader.svelte';
	import PageIndex from '$lib/layout/pageIndex.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
	import { Building2, Layers, FolderCheck, ScanEye } from 'lucide-svelte';
	import type { IndexItem } from '$lib/layout/pageIndex.svelte';
	import InfoBox from '$lib/layout/infoBox.svelte';
	import { currentUser, pb } from '@/pocketbase/index.js';
	import * as Sheet from '@/components/ui/sheet';
	import { CollectionForm } from '@/collections-components/index.js';
	import A from '@/components/ui-custom/a.svelte';

	let { data } = $props();
	const { provider, hasClaim, isClaimed } = $derived(data);

	//

	const sections: Record<string, IndexItem> = {
		general_info: {
			icon: Building2,
			anchor: 'general_info',
			label: m.General_info()
		},
		credentials: {
			icon: Layers,
			anchor: 'credentials',
			label: m.OID4VCI_Meta_Data()
		},
		test_results: {
			icon: FolderCheck,
			anchor: 'test_results',
			label: m.Test_results()
		},
		compatible_verifiers: {
			icon: ScanEye,
			anchor: 'compatible_verifiers',
			label: m.Compatible_verifiers()
		}
	};

	//

	let isClaimFormOpen = $state(false);
</script>

<PageTop>
	<div class="flex items-center gap-6">
		{#if provider.logo}
			{@const providerUrl = pb.files.getURL(provider, provider.logo)}
			<Avatar src={providerUrl} class="size-32 rounded-sm" hideIfLoadingError />
		{/if}

		<div class="space-y-3">
			<div class="space-y-1">
				<T class="text-sm">{m.Provider_name()}</T>
				<T tag="h1">{provider.name}</T>
			</div>

			{#if isClaimed}
				<div class="text-sm">
					<p class="text-muted-foreground">This provider is verified ✅</p>
				</div>
			{:else if hasClaim}
				<div class="text-sm">
					<p class="text-muted-foreground">
						You already have submitted a claim for this provider.
					</p>
					<p><A href="/my">Review your submission in the dashboard.</A></p>
				</div>
			{:else if $currentUser}
				<Button size="sm" onclick={() => (isClaimFormOpen = true)}>Claim provider</Button>
			{:else}
				<Button size="sm" href="/login" variant="outline">Login to claim provider</Button>
			{/if}
		</div>
	</div>
</PageTop>

<PageContent class="bg-secondary grow" contentClass="flex flex-col md:flex-row gap-12 items-start">
	<div class="sticky top-5">
		<PageIndex sections={Object.values(sections)} />
	</div>

	<div class="grow space-y-12">
		<div class="space-y-6">
			<PageHeader title={sections.general_info.label} id={sections.general_info.anchor} />
			<div class="flex gap-6">
				<InfoBox label="Legal entity">{provider.legal_entity}</InfoBox>
				<InfoBox label="Country">{provider.country}</InfoBox>
			</div>
			<InfoBox label={m.Description()}>{provider.description}</InfoBox>

			<div class="flex gap-6">
				<InfoBox label="Website">
					<A href={provider.external_website_url} target="_blank">
						{provider.external_website_url}
					</A>
				</InfoBox>
				<InfoBox label="Documentation">
					<A href={provider.documentation_url} target="_blank">
						{provider.documentation_url}
					</A>
				</InfoBox>
				<InfoBox label="Contact email">
					<A href={`mailto:${provider.contact_email}`} target="_blank">
						{provider.contact_email}
					</A>
				</InfoBox>
			</div>
		</div>

		<div>
			<PageHeader title={m.Apps()} id="apps"></PageHeader>
		</div>

		<div>
			<PageHeader title={m.Issuers()} id="issuers" />

			{#if provider.expand?.credential_issuers}
				<ul>
					{#each provider.expand.credential_issuers as issuer, index (issuer.id)}
						<li class="flex items-start gap-2">
							{@render CircledNumber(index + 1)}
							<div class="space-y-1">
								<InfoBox label={m.OpenID_Issuance_URL()}>
									<a
										href="/credential-issuers/{issuer.id}"
										class="hover:underline"
									>
										{issuer.url}
									</a>
								</InfoBox>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
	</div>
</PageContent>

{#snippet CircledNumber(index: number)}
	<div
		class="border-primary flex size-4 shrink-0 items-center justify-center rounded-full border text-sm text-slate-500"
	>
		<p>
			{index}
		</p>
	</div>
{/snippet}

<!-- Provider claim -->

<Sheet.Root bind:open={isClaimFormOpen}>
	<Sheet.Content class="overflow-auto sm:max-w-5xl">
		<Sheet.Header class="mb-8">
			<Sheet.Title>Claim provider</Sheet.Title>
			<Sheet.Description>
				Please fill in the following details to claim this provider.<br />
				All fields are required.
			</Sheet.Description>
		</Sheet.Header>

		<CollectionForm
			collection="provider_claims"
			onSuccess={() => {
				isClaimFormOpen = false;
			}}
			fieldsOptions={{
				order: ['name', 'description', 'logo', 'legal_entity', 'country'],
				placeholders: {
					name: 'Provider A',
					legal_entity: 'Example Org'
				},
				labels: {
					legal_entity: 'Legal entity'
				},
				exclude: ['provider', 'owner', 'status'],
				hide: {
					provider: provider.id
				}
			}}
			uiOptions={{
				hideRequiredIndicator: true,
				showToastOnSuccess: true,
				toastText: 'Provider claim request sent. Review your submission in your dashboard.'
			}}
			submitButtonContent={SubmitButton}
		/>
	</Sheet.Content>
</Sheet.Root>

{#snippet SubmitButton()}
	Send request
{/snippet}
