<script lang="ts">
	import InfoBox from '$lib/layout/infoBox.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageHeader from '$lib/layout/pageHeader.svelte';
	import PageIndex, { type IndexItem } from '$lib/layout/pageIndex.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';

	import { Building2 } from 'lucide-svelte';

	let { data } = $props();

	const sections = {
		issuer_properties: {
			label: 'Issuer properties',
			anchor: 'issuer-properties',
			icon: Building2
		}
	} satisfies Record<string, IndexItem>;
</script>

<PageTop>
	<T tag="h1">Credential issuer – {data.issuer.id}</T>
</PageTop>

<PageContent class="grow bg-secondary" contentClass="flex gap-12 items-start">
	<PageIndex sections={Object.values(sections)} class="sticky top-5" />

	<div class="grow space-y-16">
		<div class="space-y-6">
			<PageHeader
				title={sections.issuer_properties.label}
				id={sections.issuer_properties.anchor}
			/>

			<InfoBox label="URL">
				<a href={data.issuer.url} class="hover:underline" target="_blank">
					{data.issuer.url}
				</a>
			</InfoBox>
		</div>
	</div>
</PageContent>
