import type { Simplify } from 'type-fest';
import type { ListResult } from 'pocketbase';
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

/* Query runners */

export interface PocketbaseQueryAgent<
	C extends CollectionName,
	E extends PocketbaseQueryExpandOption<C> = never
> {
	getList(page: number): Promise<ListResult<PocketbaseQueryResponse<C, E>>>;
	getFullList(): Promise<PocketbaseQueryResponse<C, E>[]>;
	getOne(id: string): Promise<PocketbaseQueryResponse<C, E>>;
}

export type PocketbaseQueryAgentOptions = {
	pocketbase?: Pocketbase;
} & PocketbaseListOptions;

export function createPocketbaseQueryAgent<
	C extends CollectionName,
	E extends PocketbaseQueryExpandOption<C> = never
>(
	query: PocketbaseQuery<C, E>,
	options: PocketbaseQueryAgentOptions = {}
): PocketbaseQueryAgent<C, E> {
	const { collection } = query;

	const pbOptions: PocketbaseListOptions = {
		options,
		...buildPocketbaseQuery(query)
	};

	const pocketbase = options.pocketbase ?? pb;

	return {
		getOne: (id: string) => pocketbase.collection(collection).getOne(id, pbOptions),
		getFullList: () => pocketbase.collection(collection).getFullList(pbOptions),
		getList: (page: number) =>
			pocketbase.collection(collection).getList(page, pbOptions.perPage, pbOptions)
	};
}
