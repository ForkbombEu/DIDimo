<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import type { IconComponent } from '@/components/types.js';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
	import { Building2, Layers, FolderCheck, ScanEye } from 'lucide-svelte';

	//

	let { data } = $props();
	const { provider } = $derived(data);

	//

	const sections = {
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
	} satisfies Record<string, { icon: IconComponent; anchor: string; label: string }>;
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
		<Button size="sm">{m.Claim_service()}</Button>
	</div>
</PageTop>

<PageContent class="bg-secondary grow" contentClass="flex gap-12">
	<div>
		<ul class="space-y-4">
			{#each Object.values(sections) as section}
				<li>
					<a href={`#${section.anchor}`} class="flex items-center gap-2 hover:underline">
						<section.icon class="size-4" />
						{section.label}
					</a>
				</li>
			{/each}
		</ul>
	</div>

	<div class="grow space-y-12">
		<div>
			{@render header(sections.general_info.label, sections.general_info.anchor)}
		</div>

		<div>
			{@render header(m.Apps(), 'apps')}
		</div>

		<div>
			{@render header(m.Issuers(), 'issuers')}
			{#if provider.expand?.credential_issuers}
				<ul>
					{#each provider.expand.credential_issuers as issuer, index (issuer.id)}
						<li class="flex items-start gap-2">
							{@render CircledNumber(index + 1)}
							<div class="space-y-1">
								<T class="block text-sm">{m.OpenID_Issuance_URL()}</T>
								<a href={issuer.url} class="info block" target="_blank">
									{issuer.url}
								</a>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
	</div>
</PageContent>

{#snippet header(title: string, id: string)}
	<div {id} class="border-secondary-foreground mb-6 scroll-mt-5 border-b">
		<T tag="h2">{title}:</T>
	</div>
{/snippet}

{#snippet CircledNumber(index: number)}
	<div
		class="border-primary flex size-4 items-center justify-center rounded-full border text-sm text-slate-500"
	>
		<p>
			{index}
		</p>
	</div>
{/snippet}

<style lang="postcss">
	.info {
		@apply w-fit rounded-sm border border-slate-400 bg-white px-2 py-1;
	}
</style>
