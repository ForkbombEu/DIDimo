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
	import { currentUser } from '@/pocketbase/index.js';

	let { data } = $props();
	const { provider } = $derived(data);

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
</script>

<PageTop>
	<div class="flex items-start gap-4"></div>
	{#if provider.logo}
		<Avatar src={provider.logo} class="size-32" hideIfLoadingError />
	{/if}
	<div class="space-y-3">
		<div class="space-y-1">
			<T class="text-sm">{m.Provider_name()}</T>
			<T tag="h1">{provider.name}</T>
		</div>
		{#if $currentUser}
			<Button size="sm">Claim provider</Button>
		{:else}
			<Button size="sm" href="/login" variant="outline">Login to claim provider</Button>
		{/if}
	</div>
</PageTop>

<PageContent class="bg-secondary grow" contentClass="flex gap-12">
	<div>
		<PageIndex sections={Object.values(sections)} />
	</div>

	<div class="grow space-y-12">
		<div>
			<PageHeader title={sections.general_info.label} id={sections.general_info.anchor} />
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
