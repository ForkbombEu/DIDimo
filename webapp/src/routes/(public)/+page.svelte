<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import { CollectionManager } from '@/collections-components';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
</script>

<PageTop>
	<div class="space-y-2">
		<T tag="h1" class="text-balance">{m.Find_and_test_identity_solutions_with_ease()}</T>
		<T tag="h3" class="text-balance">
			{m.Didimo_is_your_trusted_source_for_compliance_verification()}
		</T>
	</div>
	<div class="flex gap-4">
		<Button variant="default" href="/tests/new">{m.Start_a_new_test()}</Button>
		<Button variant="secondary">{m.See_how_it_works()}</Button>
	</div>
</PageTop>

<PageContent class="bg-secondary" contentClass="space-y-12">
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_solutions()}</T>
			<Button variant="default" href="/services">{m.All_solutions()}</Button>
		</div>
		<PageGrid>
			{@const MAX_ITEMS = 3}
			<CollectionManager
				collection="services"
				queryOptions={{ perPage: MAX_ITEMS }}
				hide={['pagination']}
			>
				{#snippet records({ records })}
					{#each records as service, i}
						{@const isLast = i == MAX_ITEMS - 1}
						<ServiceCard {service} class={isLast ? 'hidden lg:block' : ''} />
					{/each}
				{/snippet}
			</CollectionManager>
		</PageGrid>
	</div>

	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_credentials()}</T>
			<Button variant="default" href="/">{m.All_credentials()}</Button>
		</div>
		<PageGrid>
			<CredentialCard class="grow basis-1" />
			<CredentialCard class="grow basis-1" />
			<CredentialCard class="hidden grow basis-1 lg:block" />
		</PageGrid>
	</div>

	<div class="space-y-6">
		<div>
			<T tag="h3">{m.Compare_by_test_results()}</T>
		</div>
		<div class="bg-card border-primary h-96 w-full rounded-lg border p-6">
			<p>table here</p>
		</div>
	</div>
</PageContent>
