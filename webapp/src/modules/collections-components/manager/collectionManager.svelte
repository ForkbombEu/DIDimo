<script
	lang="ts"
	generics="C extends CollectionName, E extends PocketbaseQueryExpandOption<C> = never"
>
	// Logic
	import { CollectionManager } from './collectionManager.svelte.js';
	import { setupComponentPocketbaseSubscriptions } from '@/pocketbase/subscriptions';
	import {
		setCollectionManagerContext,
		type CollectionManagerContext,
		type Filter,
		type FilterGroup
	} from './collectionManagerContext';

	// Logic - Types
	import type {
		PocketbaseQueryExpandOption,
		PocketbaseQueryOptions,
		PocketbaseQueryResponse
	} from '@/pocketbase/query';
	import type { CollectionName } from '@/pocketbase/collections-models';
	import type {
		UIOptions as CollectionFormUIOptions,
		FieldsOptions
	} from '../form/collectionFormTypes';
	import type { FormOptions as SuperformsOptions } from '@/forms';
	import type { CollectionFormData } from '@/pocketbase/types';

	// Components
	import Card from './recordCard.svelte';
	import Table from './table/collectionTable.svelte';
	import EmptyState from '@/components/ui-custom/emptyState.svelte';
	import Pagination from './collectionManagerPagination.svelte';
	import Search from './collectionManagerSearch.svelte';
	import Header from './collectionManagerHeader.svelte';
	import Filters from './collectionManagerFilters.svelte';

	// UI
	import { m } from '@/i18n';
	import { FolderIcon, SearchIcon, MessageCircleWarning } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	//

	type Props = {
		collection: C;
	} & Partial<Options> &
		Partial<Snippets>;

	type Snippets = {
		records: Snippet<
			[
				{
					records: PocketbaseQueryResponse<C, E>[];
					Card: typeof Card;
					Table: typeof Table;
					Pagination: typeof Pagination;
				}
			]
		>;
		emptyState: Snippet<[{ EmptyState: typeof EmptyState }]>;
		top: Snippet<[{ Search: typeof Search; Header: typeof Header; Filters: typeof Filters }]>;
		contentWrapper: Snippet<[children: () => ReturnType<Snippet>]>;
	};

	type Options = {
		queryOptions: PocketbaseQueryOptions<C, E>;
		subscribe: 'off' | 'expanded_collections' | CollectionName[];
		hide: ('empty_state' | 'pagination')[];

		formUIOptions: CollectionFormUIOptions;
		createFormUIOptions: CollectionFormUIOptions;
		editFormUIOptions: CollectionFormUIOptions;

		formSuperformsOptions: SuperformsOptions<CollectionFormData[C]>;
		createFormSuperformsOptions: SuperformsOptions<CollectionFormData[C]>;
		editFormSuperformsOptions: SuperformsOptions<CollectionFormData[C]>;

		formFieldsOptions: Partial<FieldsOptions<C>>;
		createFormFieldsOptions: Partial<FieldsOptions<C>>;
		editFormFieldsOptions: Partial<FieldsOptions<C>>;

		filters: (Filter | FilterGroup)[];
	};

	//

	const {
		collection,
		queryOptions = {},
		hide = [],
		subscribe = 'expanded_collections',
		top,
		records,
		emptyState,
		contentWrapper,
		filters = [],
		...rest
	}: Props = $props();

	//

	const manager = $derived(new CollectionManager(collection, queryOptions));

	const context = $derived<CollectionManagerContext<C, E>>({
		manager,
		filters,
		formsOptions: {
			base: {
				uiOptions: rest.formUIOptions,
				superformsOptions: rest.formSuperformsOptions,
				fieldsOptions: rest.formFieldsOptions
			},
			create: {
				uiOptions: rest.createFormUIOptions,
				superformsOptions: rest.createFormSuperformsOptions,
				fieldsOptions: rest.createFormFieldsOptions
			},
			edit: {
				uiOptions: rest.editFormUIOptions,
				superformsOptions: rest.editFormSuperformsOptions,
				fieldsOptions: rest.editFormFieldsOptions
			}
		}
	});

	setCollectionManagerContext(() => context);

	setupComponentPocketbaseSubscriptions({
		collection,
		callback: () => manager.loadRecords(),
		expandOption: queryOptions.expand,
		other: ['authorizations']
	});
</script>

{@render top?.({ Search, Header, Filters })}
{@render (contentWrapper ?? defaultContentWrapper)(content)}

{#snippet defaultContentWrapper(children: () => ReturnType<Snippet>)}
	{@render children()}
{/snippet}

{#snippet content()}
	{#if manager.loadingError}
		<EmptyState
			title={m.Error()}
			description={m.Some_error_occurred_while_loading_records_()}
			icon={MessageCircleWarning}
		/>
	{:else if manager.records.length > 0}
		{@render records?.({ records: manager.records, Card, Table, Pagination })}

		{#if !hide.includes('pagination')}
			<Pagination class="mt-6" />
		{/if}
	{:else if manager.query.hasSearch() && manager.records.length === 0}
		<EmptyState title={m.No_records_found()} icon={SearchIcon} />
	{:else if emptyState}
		{@render emptyState({ EmptyState })}
	{:else if !hide.includes('empty_state')}
		<EmptyState
			title={m.No_items_here()}
			description={m.Start_by_adding_a_record_to_this_collection_()}
			icon={FolderIcon}
		/>
	{/if}
{/snippet}
