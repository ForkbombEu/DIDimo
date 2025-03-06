import { String } from 'effect';

import type { KeyOf } from '@/utils/types';
import { getCollectionModel, type CollectionName } from '@/pocketbase/collections-models';
import type { CollectionExpands, CollectionResponses, RecordIdString } from '@/pocketbase/types';
import { ensureArray, type MaybeArray } from '@/utils/other';
import type { PocketbaseListOptions } from './utils';

/* Utils */

type Field<C extends CollectionName> = KeyOf<CollectionResponses[C]>;

/* Query */

export type PocketbaseQueryExpandOption<C extends CollectionName> = KeyOf<CollectionExpands[C]>[];

export type PocketbaseQueryOptions<
	C extends CollectionName,
	Expand extends PocketbaseQueryExpandOption<C> = never
> = Partial<{
	expand: Expand;
	perPage: number;
	filter: MaybeArray<string>;
	search: MaybeArray<SearchFilter<Field<C> | string>>;
	searchFields: Field<C>[];
	exclude: MaybeArray<ExcludeFilter>;
	sort: MaybeArray<SortOption<Field<C>>>;
}>;

export type PocketbaseQuery<
	C extends CollectionName,
	Expand extends PocketbaseQueryExpandOption<C> = never
> = PocketbaseQueryOptions<C, Expand> & {
	collection: C;
};

/* Editor */

export class PocketbaseQueryOptionsEditor<
	C extends CollectionName,
	Expand extends PocketbaseQueryExpandOption<C> = never
> {
	constructor(
		private options: PocketbaseQueryOptions<C, Expand>,
		private rootOptions: PocketbaseQueryOptions<C, Expand> = {}
	) {}

	getMergedOptions(): PocketbaseQueryOptions<C, Expand> {
		return {
			// @ts-expect-error Not relevant type error
			expand: [...ensureArray(this.rootOptions.expand), ...ensureArray(this.options.expand)],
			perPage: this.options.perPage ?? this.rootOptions.perPage,
			filter: [...ensureArray(this.rootOptions.filter), ...ensureArray(this.options.filter)],
			search: [...ensureArray(this.rootOptions.search), ...ensureArray(this.options.search)],
			searchFields: [
				...ensureArray(this.rootOptions.searchFields),
				...ensureArray(this.options.searchFields)
			],
			exclude: [
				...ensureArray(this.rootOptions.exclude),
				...ensureArray(this.options.exclude)
			],
			sort: [
				...ensureSortOptionArray(this.rootOptions.sort),
				...ensureSortOptionArray(this.options.sort)
			]
		};
	}

	//

	withPagination(perPage: number) {
		this.options.perPage = perPage;
		return this;
	}

	getPageSize() {
		return this.getMergedOptions().perPage;
	}

	hasPagination() {
		return Boolean(this.getPageSize());
	}

	//

	addFilter(filter: string) {
		this.options.filter = [...ensureArray(this.options.filter), filter];
		return this;
	}

	setFilters(filters: string | string[]) {
		this.options.filter = ensureArray(filters);
		return this;
	}

	removeFilter(filter: string) {
		this.options.filter = ensureArray(this.options.filter).filter((f) => f !== filter);
		return this;
	}

	hasFilter(filter: string) {
		return ensureArray(this.options.filter).includes(filter);
	}

	//

	addSearch(search: SearchFilter<Field<C>>) {
		this.options.search = [...ensureArray(this.options.search), search];
		return this;
	}

	setSearch(search: SearchFilter<Field<C>>) {
		this.options.search = ensureArray(search);
		return this;
	}

	clearSearch() {
		this.options.search = [];
		return this;
	}

	hasSearch() {
		return ensureArray(this.getMergedOptions().search).length > 0;
	}

	//

	addExclude(exclude: ExcludeFilter | ExcludeFilter[]) {
		this.options.exclude = [...ensureArray(this.options.exclude), ...ensureArray(exclude)];
		return this;
	}

	//

	addSort(field: Field<C>, order: SortOrder) {
		this.options.sort = [...ensureSortOptionArray(this.options.sort), [field, order]];
		return this;
	}

	setSort(field: Field<C>, order: SortOrder) {
		this.options.sort = [[field, order]];
		return this;
	}

	flipSort(sort: SortOption<Field<C>>) {
		const sorts = ensureSortOptionArray(this.options.sort);
		const sortToChange = sorts.find((s) => s[0] == sort[0] && s[1] == sort[1]);
		if (!sortToChange) return this;
		const index = sorts.indexOf(sortToChange);
		sorts[index] = [sortToChange[0], sortToChange[1] == 'ASC' ? 'DESC' : 'ASC'];
		this.options.sort = sorts;
		return this;
	}

	hasSort(field: Field<C> | string) {
		return ensureArray(this.getMergedOptions().sort).some((s) => s[0] == field);
	}

	getSort(field: Field<C> | string) {
		return ensureArray(this.getMergedOptions().sort).find((s) => s[0] == field);
	}
}

//

const QUOTE = "'";

//

type BaseSearchFilter<T extends string = string> = {
	text: string;
	fields: T[];
};

type SearchFilter<T extends string = string> = BaseSearchFilter<T> | string;

function buildSearchFilter(filter: BaseSearchFilter) {
	const search = filter.fields.map((f) => `${f} ~ ${QUOTE}${filter.text}${QUOTE}`).join(' || ');
	return `(${search})`;
}

//

type ExcludeFilter = RecordIdString;

function buildExcludeFilter(filter: ExcludeFilter) {
	return `id != ${QUOTE}${filter}${QUOTE}`;
}

//

type SortOrder = 'ASC' | 'DESC';

type SortOption<T extends string = string> = [T, SortOrder];

function buildSortOption<T extends string>(sortOption: SortOption<T>) {
	return (sortOption[1] == 'ASC' ? '+' : '-') + sortOption[0];
}

/* Build */

export function buildPocketbaseQuery<
	C extends CollectionName,
	Expand extends PocketbaseQueryExpandOption<C>
>(query: PocketbaseQuery<C, Expand>): PocketbaseListOptions {
	const { collection, ...options } = query;
	//

	const allCollectionFields = getCollectionModel(collection).fields.map(
		(f) => f.name
	) as Field<C>[];

	const baseSearchFields =
		query.searchFields && query.searchFields.length > 0
			? query.searchFields
			: allCollectionFields;

	const filter = [
		...ensureArray(options.filter),
		...ensureArray(options.exclude).map(buildExcludeFilter),
		...ensureArray(options.search)
			.map((searchFilter) => {
				if (typeof searchFilter == 'string')
					return {
						text: searchFilter,
						fields: baseSearchFields
					};
				else return searchFilter;
			})
			.map(buildSearchFilter)
	].join(' && ');

	//

	const sort = ensureSortOptionArray(options.sort).map(buildSortOption).join(',');

	const expand = ensureArray(options.expand).join(',');

	//

	const listOptions: PocketbaseListOptions = {};

	if (options.perPage) listOptions.perPage = options.perPage;
	if (String.isNonEmpty(expand)) listOptions.expand = expand;
	if (String.isNonEmpty(filter)) listOptions.filter = `(${filter})`;
	if (String.isNonEmpty(sort)) listOptions.sort = sort;

	return listOptions;
}

/* Utils */

function isSortOption(unknown: unknown): unknown is SortOption<string> {
	return (
		Array.isArray(unknown) &&
		unknown.length == 2 &&
		typeof unknown[0] == 'string' &&
		(unknown[1] == 'ASC' || unknown[1] == 'DESC')
	);
}

function ensureSortOptionArray<T extends string = string>(
	unknown: MaybeArray<SortOption<T>>
): SortOption<T>[] {
	if (isSortOption(unknown)) return [unknown];
	else return ensureArray(unknown);
}
