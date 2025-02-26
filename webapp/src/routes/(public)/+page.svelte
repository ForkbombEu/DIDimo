<script lang="ts">
	import CredentialCard from '$lib/layout/credentialCard.svelte';
	import FakeTable from '$lib/layout/fakeTable.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageGrid from '$lib/layout/pageGrid.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import ServiceCard from '$lib/layout/serviceCard.svelte';
	import Alert from '@/components/ui-custom/alert.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { featureFlags } from '@/features';
	import { createForm, Form, SubmitButton } from '@/forms';
	import { Field } from '@/forms/fields';
	import { m } from '@/i18n';
	import LanguageSelect from '@/i18n/languageSelect.svelte';
	import { pb } from '@/pocketbase';
	import { Collections, ServicesCountryOptions, type ServicesResponse } from '@/pocketbase/types';
	import { onMount } from 'svelte';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';

	const fakeService: ServicesResponse = {
		id: 'das',
		country: ServicesCountryOptions.IT,
		created: '2024-12-12',
		updated: '2024-12-12',
		credential_issuers: [],
		description: 'Lorem ipsum',
		external_links: [],
		legal_entity: 'ForkbombEu',
		logo: 'https://avatars.githubusercontent.com/u/96812851?s=200&v=4',
		name: 'Test credential issuer',
		owner: 'id',
		collectionId: '',
		collectionName: Collections.Services
	};

	const schema = z.object({
		name: z.string(),
		email: z.string().email()
	});

	const form = createForm({
		adapter: zod(schema),
		onSubmit: async ({ form: { data } }) => {
			try {
				await pb.collection('waitlist').create({
					email: data.email,
					name: data.name
				});
				formSuccess = true;
			} catch {
				throw new Error(
					m.An_error_occurred_while_submitting_your_request_Please_try_again()
				);
			}
		}
	});

	let formSuccess = $state(false);

	//

	let formHighlight = $state(false);

	onMount(() => {
		// TODO (start animation on scroll)
	});
</script>

{#if !$featureFlags.DEMO}
	<div class="flex justify-end px-6 pt-4">
		<LanguageSelect flagsOnly>
			{#snippet trigger({ triggerAttributes, language })}
				<Button variant="outline" class="w-14 text-2xl" {...triggerAttributes}>
					{language.flag}
				</Button>
			{/snippet}
		</LanguageSelect>
	</div>
{/if}
<PageTop>
	<div class="space-y-2">
		<T tag="h1" class="text-balance">{m.Find_and_test_identity_solutions_with_ease()}</T>
		<T tag="h3" class="text-balance">
			{m.Didimo_is_your_trusted_source_for_compliance_verification()}
		</T>
	</div>
	<div class="flex gap-4">
		<Button variant="default" href={$featureFlags.DEMO ? '#waitlist' : '/tests/new'}
			>{m.Start_a_new_test()}</Button
		>
		<Button variant="secondary">{m.See_how_it_works()}</Button>
	</div>
</PageTop>

<PageContent class="bg-secondary" contentClass="space-y-12">
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_solutions()}</T>

			{#if $featureFlags.DEMO}
				<Button variant="default" disabled class="select-none blur">
					{m.All_solutions()}
				</Button>
			{:else}
				<Button variant="default" href="/providers">
					{m.All_solutions()}
				</Button>
			{/if}
		</div>

		<PageGrid class="select-none blur-sm">
			<ServiceCard service={fakeService} class="grow basis-1" />
			<ServiceCard service={fakeService} class="grow basis-1" />
			<ServiceCard service={fakeService} class="hidden grow basis-1 lg:block" />
			<!--
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
			</CollectionManager> -->
		</PageGrid>
	</div>

	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<T tag="h3">{m.Find_credentials()}</T>
			{#if $featureFlags.DEMO}
				<Button variant="default" disabled class="select-none blur">
					{m.All_credentials()}
				</Button>
			{:else}
				<Button variant="default" href="/credentials">{m.All_credentials()}</Button>
			{/if}
		</div>
		<PageGrid>
			{@const MAX_ITEMS = 3}
			<CollectionManager
				collection="credentials"
				queryOptions={{ perPage: MAX_ITEMS }}
				hide={['pagination']}
			>
				{#snippet records({ records })}
					{#each records as credential, i}
						{@const isLast = i == MAX_ITEMS - 1}
						<CredentialCard {credential} class={isLast ? 'hidden lg:flex' : ''} />
					{/each}
				{/snippet}
			</CollectionManager>
		</PageGrid>
	</div>
</PageContent>

<PageContent class="border-y-primaryborder-y-2" contentClass="!space-y-8">
	<div id="waitlist" class="scroll-mt-20">
		<T tag="h2" class="text-balance">
			{m._Stay_Ahead_in_Digital_Identity_Compliance_Join_Our_Early_Access_List()}
		</T>
		<T class="mt-1 text-balance font-medium">
			{m.Be_the_first_to_explore_DIDimo_the_ultimate_compliance_testing_tool_for_decentralized_identity_Get_exclusive_updates_early_access_and_a_direct_line_to_our_team_()}
		</T>
	</div>

	{#if !formSuccess}
		<Form {form} hide={['submit_button']} class=" !space-y-3" hideRequiredIndicator>
			<div class="flex w-full max-w-3xl flex-col gap-2 md:flex-row md:gap-6">
				<div class="grow">
					<Field
						{form}
						name="name"
						options={{
							label: m.Your_name(),
							placeholder: m.John_Doe(),
							class: 'bg-secondary/40 '
						}}
					/>
				</div>
				<div class="grow">
					<Field
						{form}
						name="email"
						options={{
							label: m.Your_email(),
							placeholder: m.e_g_hellomycompany_com(),
							class: 'bg-secondary/40'
						}}
					/>
				</div>
			</div>
			<SubmitButton>{m.Join_the_Waitlist()}</SubmitButton>
		</Form>
	{:else}
		<Alert variant="info">
			<p class="font-bold">{m.Request_sent_()}</p>
			<p>
				{m.Thanks_for_your_interest_We_will_write_to_you_soon()}
			</p>
		</Alert>
	{/if}
</PageContent>

<PageContent class="bg-secondary" contentClass="space-y-12">
	<div class="space-y-6">
		<div>
			<T tag="h3">{m.Compare_by_test_results()}</T>
		</div>
		<FakeTable />
	</div>
</PageContent>
