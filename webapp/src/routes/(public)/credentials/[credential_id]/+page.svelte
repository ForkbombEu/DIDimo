<script lang="ts">
	import InfoBox from '$lib/layout/infoBox.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageHeader from '$lib/layout/pageHeader.svelte';
	import PageIndex from '$lib/layout/pageIndex.svelte';
	import type { IndexItem } from '$lib/layout/pageIndex.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import type { CredentialConfiguration } from '$lib/types/openid.js';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m, localizeHref } from '@/i18n/index.js';
	import { Building2, FolderCheck, Layers3, ScanEye } from 'lucide-svelte';

	let { data } = $props();
	const { credential } = $derived(data);

	const sections = {
		credential_properties: {
			icon: Building2,
			anchor: 'credential_properties',
			label: 'Credential properties'
		},
		credential_subjects: {
			icon: Layers3,
			anchor: 'credential_subject',
			label: 'Credential subject'
		},
		// compatible_apps: {
		// 	icon: Building2,
		// 	anchor: 'compatible_apps',
		// 	label: 'Compatible apps'
		// },
		compatible_issuer: {
			icon: FolderCheck,
			anchor: 'compatible_issuer',
			label: 'Compatible issuer'
		}
		// test_results: {
		// 	icon: ScanEye,
		// 	anchor: 'test_results',
		// 	label: 'Test results'
		// }
	} satisfies Record<string, IndexItem>;

	const credentialConfiguration = $derived(
		credential.json as CredentialConfiguration | undefined
	);
	const credentialSubject = $derived(
		credentialConfiguration?.credential_definition?.credentialSubject
	);
	const credentialIssuer = $derived(credential.expand?.credential_issuer);
</script>

<PageTop>
	<div class="flex items-center gap-2">
		<Avatar src={credential.logo} class="!rounded-sm" hideIfLoadingError />
		<div class="">
			<T class="">Credential name</T>
			<T tag="h1">{credential.name}</T>
		</div>
	</div>
</PageTop>

<PageContent class="grow bg-secondary" contentClass="flex gap-12 items-start">
	<PageIndex sections={Object.values(sections)} class="sticky top-5" />

	<div class="grow space-y-16">
		<div class="space-y-6">
			<PageHeader
				title={sections.credential_properties.label}
				id={sections.credential_properties.anchor}
			/>

			<div class="flex gap-6">
				<InfoBox label="Issuer" value={credential.issuer_name} />
				<InfoBox label="Format" value={credential.format} />
				<InfoBox label="Locale" value={credential.locale} />
			</div>

			<div class="flex gap-6">
				<InfoBox
					label="Signing algorithms supported"
					value={credentialConfiguration?.credential_signing_alg_values_supported?.join(
						', '
					)}
				/>
				<InfoBox
					label="Cryptographic binding methods supported"
					value={credentialConfiguration?.cryptographic_binding_methods_supported?.join(
						', '
					)}
				/>
			</div>

			<InfoBox label="Description" value={credential.description} />
			<InfoBox label="Type" value={credential.type} />
		</div>

		<div class="space-y-6">
			<PageHeader
				title={sections.credential_subjects.label}
				id={sections.credential_subjects.anchor}
			/>

			{#if credentialSubject}
				<InfoBox
					label="Type"
					value={credentialConfiguration?.credential_definition?.type?.join(', ')}
				/>

				<div class="grid grid-cols-[auto_auto_auto] gap-3">
					{#each Object.entries(credentialSubject) as [key, value]}
						<InfoBox label="Property">
							<T>{key}</T>
						</InfoBox>

						{#if value.display}
							<InfoBox label="Label">
								<T>
									{value.display
										.map((d) => `${d.name}${d.locale ? ` (${d.locale})` : ''}`)
										.join(', ')}
								</T>
							</InfoBox>
						{:else}
							<div></div>
						{/if}

						{#if value.mandatory}
							<InfoBox label="Required">
								<T>Mandatory</T>
							</InfoBox>
						{:else}
							<div></div>
						{/if}
					{/each}
				</div>
			{/if}
		</div>

		<!-- <div>
			<PageHeader
				title={sections.compatible_apps.label}
				id={sections.compatible_apps.anchor}
			/>
		</div> -->

		<div>
			<PageHeader
				title={sections.compatible_issuer.label}
				id={sections.compatible_issuer.anchor}
			/>

			{#if credentialIssuer}
				<InfoBox label="OpenID issuance URL">
					<a
						href={localizeHref(`/credential-issuers/${credentialIssuer.id}`)}
						class="hover:underline"
					>
						{credentialIssuer.url}
					</a>
				</InfoBox>
			{/if}
		</div>

		<!-- <div>
			<PageHeader title={sections.test_results.label} id={sections.test_results.anchor} />
		</div> -->
	</div>
</PageContent>
