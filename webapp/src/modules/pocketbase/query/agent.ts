import type { Simplify } from 'type-fest';
import type Pocketbase from 'pocketbase';
import type { CollectionResponses, CollectionExpands } from '@/pocketbase/types';
import {
	buildPocketbaseQuery,
	type PocketbaseQuery,
	type PocketbaseQueryExpandOption
} from './query';
import type { CollectionName } from '@/pocketbase/collections-models';
import { pb } from '@/pocketbase';
import type { PocketbaseListOptions } from './utils';

/* Query response */

export type PocketbaseQueryResponse<
	C extends CollectionName,
	E extends PocketbaseQueryExpandOption<C> = never
> = CollectionResponses[C] &
	Simplify<{
		expand?: Partial<Pick<CollectionExpands[C], E[number]>>;
	}>;

/* Query agent */

export type PocketbaseQueryAgentOptions = {
	pocketbase?: Pocketbase;
} & PocketbaseListOptions;

export class PocketbaseQueryAgent<
	C extends CollectionName,
	E extends PocketbaseQueryExpandOption<C> = never
> {
	private pocketbase: Pocketbase;
	readonly collection: C;
	readonly listOptions: PocketbaseListOptions;

	constructor(query: PocketbaseQuery<C, E>, options: PocketbaseQueryAgentOptions = {}) {
		this.collection = query.collection;
		this.pocketbase = options.pocketbase ?? pb;
		this.listOptions = {
			...options,
			...buildPocketbaseQuery(query)
		};
	}

	getOne(id: string) {
		return this.pocketbase
			.collection(this.collection)
			.getOne<PocketbaseQueryResponse<C, E>>(id, this.listOptions);
	}

	getFullList() {
		return this.pocketbase
			.collection(this.collection)
			.getFullList<PocketbaseQueryResponse<C, E>>(this.listOptions);
	}

	getList(page: number) {
		return this.pocketbase
			.collection(this.collection)
			.getList<
				PocketbaseQueryResponse<C, E>
			>(page, this.listOptions.perPage, this.listOptions);
	}
}
