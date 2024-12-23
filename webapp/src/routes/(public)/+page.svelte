<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
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

<PageContent class="bg-muted" contentClass="space-y-12">
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_solutions()}</T>
			<Button variant="default" href="/">{m.All_solutions()}</Button>
		</div>
		<div class="flex flex-col gap-6 sm:flex-row md:gap-8">
			<ServiceCard />
			<ServiceCard />
			<ServiceCard class="hidden lg:block" />
		</div>
	</div>

	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_credentials()}</T>
			<Button variant="default" href="/">{m.All_credentials()}</Button>
		</div>
		<div class="flex flex-col gap-6 sm:flex-row md:gap-8">
			<CredentialCard class="grow basis-1" />
			<CredentialCard class="grow basis-1" />
			<CredentialCard class="hidden grow basis-1 lg:block" />
		</div>
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
