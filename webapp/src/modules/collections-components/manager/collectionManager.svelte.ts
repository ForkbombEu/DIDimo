import { pb } from '@/pocketbase';
import type { CollectionName } from '@/pocketbase/collections-models';
import {
	createPocketbaseQueryAgent,
	type PocketbaseQueryOptions,
	type PocketbaseQueryExpandOption,
	type PocketbaseQueryResponse,
	PocketbaseQueryOptionsEditor,
	type PocketbaseQueryAgentOptions
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

	private queryOptions: PocketbaseQueryOptions<C, E> = $state({});

	query = $derived.by(
		() => new PocketbaseQueryOptionsEditor(this.queryOptions, this.options.query)
	);
	private queryAgent = $derived.by(() =>
		createPocketbaseQueryAgent(
			{
				collection: this.collection,
				...this.query.getMergedOptions()
			},
			this.options.queryAgent
		)
	);

	constructor(
		public readonly collection: C,
		private readonly options: {
			query: PocketbaseQueryOptions<C, E>;
			queryAgent: PocketbaseQueryAgentOptions;
		}
	) {
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

	async loadRecords() {
		try {
			if (this.query.hasPagination()) {
				const result = await this.queryAgent.getList(this.currentPage);
				this.totalItems = result.totalItems;
				this.records = result.items;
			} else {
				this.records = await this.queryAgent.getFullList();
			}
		} catch (e) {
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
