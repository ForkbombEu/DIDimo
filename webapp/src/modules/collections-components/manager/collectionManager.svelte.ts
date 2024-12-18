import { pb } from '@/pocketbase';
import type { CollectionName } from '@/pocketbase/collections-models';
import {
	PocketbaseQuery,
	type ExpandQueryOption,
	type PocketbaseQueryOptions,
	type QueryResponse
} from '@/pocketbase/query';
import type { CollectionResponses, RecordIdString } from '@/pocketbase/types';
import type { ClientResponseError, RecordService } from 'pocketbase';
import { Array } from 'effect';

//

export class CollectionManager<
	C extends CollectionName,
	Expand extends ExpandQueryOption<C> = never
> {
	collection: C;
	recordService: RecordService<CollectionResponses[C]>;

	records = $state<QueryResponse<C, Expand>[]>([]);
	currentPage = $state(1);
	totalItems = $state(0);
	loadingError = $state<ClientResponseError>();

	queryOptions = $state<Partial<PocketbaseQueryOptions<C, Expand>>>({});
	query = $derived.by(() => {
		return new PocketbaseQuery(this.collection, this.queryOptions);
	});

	selectedRecords = $state<RecordIdString[]>([]);

	//

	constructor(collection: C, queryOptions: Partial<PocketbaseQueryOptions<C, Expand>> = {}) {
		this.collection = collection;
		this.recordService = pb.collection(collection);
		this.queryOptions = queryOptions;

		$effect(() => {
			this.loadRecords();
		});
	}

	async loadRecords() {
		const hasPagination = Boolean(this.queryOptions.perPage);
		try {
			if (hasPagination) {
				const result = await this.query.getList(this.currentPage);
				this.totalItems = result.totalItems;
				this.records = result.items;
			} else {
				this.records = await this.query.getFullList();
			}
		} catch (e) {
			this.loadingError = e as ClientResponseError;
		}
	}

	/* Search */

	search(text: string) {
		this.queryOptions.search = text;
	}

	clearSearch() {
		this.queryOptions.search = undefined;
	}

	/* Selection */

	areAllRecordsSelected() {
		return this.records.every((r) => this.selectedRecords.includes(r.id));
	}

	toggleSelectAllRecords() {
		const allSelected = this.areAllRecordsSelected();
		if (allSelected) {
			this.selectedRecords = [];
		} else {
			this.selectedRecords = this.records.map((r) => r.id);
		}
	}

	discardSelection() {
		this.selectedRecords = [];
	}

	selectRecord(id: RecordIdString) {
		this.selectedRecords.push(id);
	}

	deselectRecord(id: RecordIdString) {
		this.selectedRecords = Array.remove(this.selectedRecords, this.selectedRecords.indexOf(id));
	}
}
