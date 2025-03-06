import { pb } from '@/pocketbase';
import type { CollectionName } from '@/pocketbase/collections-models';
import {
	type PocketbaseQueryOptions,
	type PocketbaseQueryExpandOption,
	type PocketbaseQueryResponse,
	PocketbaseQueryOptionsEditor,
	type PocketbaseQueryAgentOptions,
	PocketbaseQueryAgent
} from '@/pocketbase/query';
import type { RecordIdString } from '@/pocketbase/types';
import type { ClientResponseError, RecordService } from 'pocketbase';
import { Array } from 'effect';

//

export class CollectionManager<
	C extends CollectionName,
	E extends PocketbaseQueryExpandOption<C> = never
> {
	recordService: RecordService<PocketbaseQueryResponse<C, E>>;

	private rootQueryOptions: PocketbaseQueryOptions<C, E> = $state({});
	private currentQueryOptions: PocketbaseQueryOptions<C, E> = $state({});
	query = $derived.by(
		() => new PocketbaseQueryOptionsEditor(this.currentQueryOptions, this.rootQueryOptions)
	);

	private queryAgentOptions: PocketbaseQueryAgentOptions = $state({});
	private queryAgent = $derived.by(
		() =>
			new PocketbaseQueryAgent(
				{
					collection: this.collection,
					...this.query.getMergedOptions()
				},
				{ ...this.queryAgentOptions, requestKey: null }
			)
	);

	constructor(
		public readonly collection: C,
		options: {
			query: PocketbaseQueryOptions<C, E>;
			queryAgent: PocketbaseQueryAgentOptions;
		}
	) {
		this.rootQueryOptions = options.query;
		this.queryAgentOptions = options.queryAgent;
		this.recordService = pb.collection(collection);

		$effect(() => {
			this.loadRecords();
		});
	}

	/* Data loading */

	records = $state<PocketbaseQueryResponse<C, E>[]>([]);
	currentPage = $state(1);
	totalItems = $state(0);
	loadingError = $state<ClientResponseError>();

	private previousFilter: string | undefined;

	async loadRecords() {
		const currentFilter = this.queryAgent.listOptions.filter;
		if (this.previousFilter !== currentFilter) {
			this.currentPage = 1;
			this.previousFilter = currentFilter;
		}

		try {
			if (this.query.hasPagination()) {
				const result = await this.queryAgent.getList(this.currentPage);
				this.totalItems = result.totalItems;
				this.records = result.items;
			} else {
				this.records = await this.queryAgent.getFullList();
			}
		} catch (e) {
			console.error(e);
			this.loadingError = e as ClientResponseError;
		}
	}

	/* Selection */

	selectedRecords = $state<RecordIdString[]>([]);

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
