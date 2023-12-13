<script lang="ts" context="module">
	import { getContext } from 'svelte';
	import type { RecordService } from 'pocketbase';
	import type { Writable } from 'svelte/store';
	import type { FieldsSettings } from '$lib/recordForm';
	import type { RecordFullListOptions } from 'pocketbase';

	export const RECORDS_MANAGER_KEY = Symbol('rmk');

	export type RecordsManagerContext<R extends PBResponse = PBResponse> = {
		collection: string;
		dataManager: {
			recordService: RecordService;
			loadRecords: () => Promise<void>;
			perPage: Writable<number>;
			currentPage: Writable<string>;
			totalPages: Writable<number>;
			totalItems: Writable<number>;
			queryParams: Writable<RecordFullListOptions>;
		};
		selectionManager: {
			selectedRecords: Writable<string[]>;
			allRecordsSelected: (selectedRecords: string[]) => boolean;
			toggleSelectAllRecords: () => void;
			discardSelection: () => void;
		};
		formFieldsSettings: {
			base: Partial<FieldsSettings<R>>;
			create: Partial<FieldsSettings<R>>;
			edit: Partial<FieldsSettings<R>>;
		};
	};

	export function getRecordsManagerContext<
		R extends PBResponse = PBResponse
	>(): RecordsManagerContext<R> {
		return getContext(RECORDS_MANAGER_KEY);
	}
</script>

<script lang="ts">
	import { onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { browser } from '$app/environment';
	import { page } from '$app/stores';

	import { pb } from '$lib/pocketbase';
	import type { Collections } from '$lib/pocketbase/types';
	import type { PBResponse } from '$lib/utils/types';
	import { createTypeProp } from '$lib/utils/typeProp';

	import GridSpinner from '$lib/components/gridSpinner.svelte';
	import CollectionEmptyState from './ui/collectionEmptyState.svelte';

	import { ExclamationTriangle } from 'svelte-heros-v2';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	//

	export let collection: Collections | string;

	export let formSettings: Partial<FieldsSettings<RecordGeneric>> = {};
	export let createFormSettings: Partial<FieldsSettings<RecordGeneric>> = {};
	export let editFormSettings: Partial<FieldsSettings<RecordGeneric>> = {};

	export let initialQueryParams: RecordFullListOptions = {};
	export let subscribe: string[] = [];

	export let perPage = 5;
	export let disablePagination = false;

	/* Data load */

	let currentPage = writable('');
	let totalItems = writable(0);
	$: currentPage.set($page.url.searchParams.get('page') || '1');

	const queryParams = writable<RecordFullListOptions>({
		sort: '-created',
		...initialQueryParams
	});

	$: $queryParams = {
		$autoCancel: false,
		...$queryParams,
		...initialQueryParams
	};

	const recordService = pb.collection(collection);

	let records: RecordGeneric[] = [];
	let totalPages = writable(0);

	async function loadRecords() {
		if (!disablePagination) {
			const res = await recordService.getList<RecordGeneric>(Number($currentPage), perPage, {
				...$queryParams
			});
			records = res.items;
			totalPages.set(res.totalPages);
			totalItems.set(res.totalItems);
		} else {
			const res = await recordService.getFullList<RecordGeneric>({
				...$queryParams
			});
			records = res;
		}
	}

	let promise = loadRecords();

	$: if (browser) {
		$queryParams;
		initialQueryParams;
		$currentPage;
		$totalPages;
		$totalItems;
		loadRecords();
	}

	onMount(() => {
		const collections = [...subscribe, collection];
		for (const c of collections) {
			pb.collection(c).subscribe('*', () => {
				loadRecords();
			});
		}

		return () => {
			for (const c of collections) {
				pb.collection(c).unsubscribe();
			}
		};
	});

	/* Record selection */

	const selectedRecords = writable<string[]>([]);

	function allRecordsSelected(selectedRecords: string[] = []) {
		return records.every((r) => selectedRecords.includes(r.id));
	}

	function toggleSelectAllRecords() {
		const allSelected = allRecordsSelected($selectedRecords);
		if (allSelected) {
			$selectedRecords = [];
		} else {
			$selectedRecords = records.map((item) => item.id);
		}
	}

	function discardSelection() {
		$selectedRecords = [];
	}

	/* Context */

	setContext<RecordsManagerContext<RecordGeneric>>(RECORDS_MANAGER_KEY, {
		collection,
		dataManager: {
			recordService,
			loadRecords,
			queryParams,
			perPage: writable(perPage),
			currentPage,
			totalPages,
			totalItems
		},
		selectionManager: {
			selectedRecords,
			allRecordsSelected,
			toggleSelectAllRecords,
			discardSelection
		},
		formFieldsSettings: {
			base: formSettings,
			create: createFormSettings,
			edit: editFormSettings
		}
	});

	//
</script>

{#await promise}
	<div class="flex items-center justify-center min-h-screen justify-items-center">
		<GridSpinner />
	</div>
{:then}
	<slot {records} {loadRecords} />
{:catch}
	<CollectionEmptyState
		icon={ExclamationTriangle}
		title="Loading error"
		description="Some error occurred while loading records."
	/>
{/await}
