<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageHeader from '$lib/layout/pageHeader.svelte';
	import PageIndex from '$lib/layout/pageIndex.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m, localizeHref } from '@/i18n';
	import { Building2, Layers, FolderCheck, ScanEye } from 'lucide-svelte';
	import type { IndexItem } from '$lib/layout/pageIndex.svelte';
	import InfoBox from '$lib/layout/infoBox.svelte';
	import { pb } from '@/pocketbase/index.js';
	import type { OrganizationInfoResponse } from '@/pocketbase/types';

	type Props = {
		organizationInfo: OrganizationInfoResponse;
	};

	let { organizationInfo }: Props = $props();

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

	const credentialIssuersPromise = pb.collection('credential_issuers').getFullList({
		filter: `owner = '${organizationInfo.owner}'`
	});
</script>

<PageTop>
	<div class="flex items-center gap-6">
		{#if organizationInfo.logo}
			{@const providerUrl = pb.files.getURL(organizationInfo, organizationInfo.logo)}
			<Avatar src={providerUrl} class="size-32 rounded-sm" hideIfLoadingError />
		{/if}

		<div class="space-y-3">
			<div class="space-y-1">
				<T class="text-sm">{m.Organization_name()}</T>
				<T tag="h1">{organizationInfo.name}</T>
			</div>
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
				<InfoBox label="Legal entity">{organizationInfo.legal_entity}</InfoBox>
				<InfoBox label="Country">{organizationInfo.country}</InfoBox>
			</div>
			<InfoBox label={m.Description()}>{organizationInfo.description}</InfoBox>

			<div class="flex gap-6">
				<InfoBox label="Website">
					<a href={organizationInfo.external_website_url} target="_blank">
						{organizationInfo.external_website_url}
					</a>
				</InfoBox>
				<InfoBox label="Contact email">
					<a href={`mailto:${organizationInfo.contact_email}`} target="_blank">
						{organizationInfo.contact_email}
					</a>
				</InfoBox>
			</div>
		</div>

		<div>
			<PageHeader title={m.Apps()} id="apps"></PageHeader>
		</div>

		<div>
			<PageHeader title={m.Issuers()} id="issuers" />

			{#await credentialIssuersPromise then credential_issuers}
				<ul>
					{#each credential_issuers as issuer, index (issuer.id)}
						<li class="flex items-start gap-2">
							{@render CircledNumber(index + 1)}
							<div class="space-y-1">
								<InfoBox label={m.OpenID_Issuance_URL()}>
									<a
										href={localizeHref(`/credential-issuers/${issuer.id}`)}
										class="hover:underline"
									>
										{issuer.url}
									</a>
								</InfoBox>
							</div>
						</li>
					{/each}
				</ul>
			{/await}
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
