<script lang="ts">
	import PageTop from '$lib/layout/pageTop.svelte';
	import { m } from '@/i18n';
	import type { PageData } from './$types';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import { pb } from '@/pocketbase';
	import PageContent from '$lib/layout/pageContent.svelte';

	let { data }: { data: PageData } = $props();
	const { service } = $derived(data);

	const serviceLogo = $derived(pb.files.getURL(service, service.logo));
</script>

<PageTop>
	<div class="space-y-4">
		<div>
			<T>{m.Service_name()}</T>
			<T tag="h1">{service.name}</T>
		</div>

		<Button>{m.Claim_service()}</Button>
	</div>
</PageTop>

<PageContent class="bg-secondary grow">
	<pre>{JSON.stringify(service, null, 4)}</pre>
</PageContent>
