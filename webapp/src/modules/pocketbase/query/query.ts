import type { RecordFullListOptions, RecordListOptions } from 'pocketbase';
import type { Simplify } from 'type-fest';
import { String } from 'effect';
import type Pocketbase from 'pocketbase';

import type { KeyOf } from '@/utils/types';
import { pb } from '@/pocketbase';
import { getCollectionModel, type CollectionName } from '@/pocketbase/collections-models';
import type { CollectionExpands, CollectionResponses, RecordIdString } from '@/pocketbase/types';

//

export class PocketbaseQuery<
	C extends CollectionName,
	Expand extends CollectionExpand<C> = never,
	Field extends string = keyof CollectionResponses[C] & string
> {
	constructor(
		public readonly collection: C,
		private readonly init: Partial<{
			expand: Expand;
			perPage: number;
			filters: string[] | string;
			search: SearchFilter<Field> | BaseSearchFilter<Field>[];
			searchFields: Field[];
			exclude: ExcludeFilter | ExcludeFilter[];
			sort: SortOption<Field> | SortOption<Field>[];
			pb: Pocketbase;
			fetch: typeof fetch;
		}>
	) {}

	//

	private perPage: number | undefined = undefined;

	withPagination(perPage: number) {
		this.perPage = perPage;
		return this;
	}

	//

	private baseFilters: string[] = [];

	addFilter(filter: string) {
		this.baseFilters.push(filter);
		return this;
	}

	setFilters(filters: string | string[]) {
		this.baseFilters = Array.isArray(filters) ? filters : [filters];
		return this;
	}

	private getAllBaseFilters() {
		const allBaseFilters = [];
		if (typeof this.init.filters == 'string') allBaseFilters.push(this.init.filters);
		else if (Array.isArray(this.init.filters)) allBaseFilters.push(...this.init.filters);
		allBaseFilters.push(...this.baseFilters);
		return allBaseFilters;
	}

	//

	private searchFilters: BaseSearchFilter[] = [];

	addSearch(text: string, searchFields: Field[] | undefined = undefined) {
		const fields = searchFields ?? this.init.searchFields ?? this.getAllSearchFields();

		this.searchFilters.push({ text, fields });
		return this;
	}

	setSearch(text: string, searchFields: Field[] | undefined = undefined) {
		this.searchFilters = [];
		return this.addSearch(text, searchFields);
	}

	private getAllSearchFields() {
		return getCollectionModel(this.collection).fields.map((f) => f.name);
	}

	private getAllSearchFilters() {
		const allSearchFilters: BaseSearchFilter[] = [];
		if (this.init.search) {
			if (typeof this.init.search == 'string')
				allSearchFilters.push({
					text: this.init.search,
					fields: this.getAllSearchFields()
				});
			else if ('text' in this.init.search) allSearchFilters.push(this.init.search);
			else if (Array.isArray(this.init.search)) allSearchFilters.push(...this.init.search);
		}
		allSearchFilters.push(...this.searchFilters);
		return allSearchFilters;
	}

	//

	private excludeFilters: ExcludeFilter[][] = [];

	addExclude(exclude: ExcludeFilter[]) {
		this.excludeFilters.push(exclude);
		return this;
	}

	private getAllExcludeFilters() {
		const allExcludeFilters = [];
		if (typeof this.init.exclude == 'string') allExcludeFilters.push(this.init.exclude);
		else if (Array.isArray(this.init.exclude)) allExcludeFilters.push(...this.init.exclude);
		allExcludeFilters.push(...this.excludeFilters.flat());
		return allExcludeFilters;
	}

	//

	private sorts: SortOption[] = [];

	addSort(field: Field, order: SortOrder) {
		this.sorts.push([field, order]);
		return this;
	}

	setSort(field: Field, order: SortOrder) {
		this.sorts = [];
		return this.addSort(field, order);
	}

	flipSort(index = 0) {
		const sortToChange = this.sorts.at(index);
		if (!sortToChange) return this;
		this.sorts[index] = [sortToChange[0], sortToChange[1] == 'ASC' ? 'DESC' : 'ASC'];
		return this;
	}

	private getAllSorts() {
		const allSorts: SortOption[] = [];
		if (this.init.sort) {
			if (this.init.sort.length == 2 && this.init.sort.every((v) => !Array.isArray(v))) {
				allSorts.push(this.init.sort as SortOption);
			} else {
				allSorts.push(...(this.init.sort as SortOption[]));
			}
		}
		allSorts.push(...this.sorts);
		return allSorts;
	}

	//

	private getAllFilters() {
		const allFilters = [
			...this.getAllBaseFilters(),
			...this.getAllSearchFilters().map(buildSearchFilter),
			...this.getAllExcludeFilters().map(buildExcludeFilter)
		];
		return allFilters.filter(String.isNonEmpty).join(' && ');
	}

	buildOptions(): PocketbaseListOptions {
		const options: PocketbaseListOptions = {
			perPage: this.perPage ?? this.init.perPage
		};

		if (this.init.expand && this.init.expand.length > 0)
			options.expand = this.init.expand.join(',');
		if (this.getAllSorts().length > 0)
			options.sort = this.getAllSorts().map(buildSortOption).join(',');
		if (String.isNonEmpty(this.getAllFilters())) options.filter = this.getAllFilters();

		return options;
	}

	//

	get pb() {
		return this.init.pb ?? pb;
	}

	getList(page: number = 0) {
		return this.pb
			.collection(this.collection)
			.getList<
				QueryResponse<C, Expand>
			>(page, this.perPage ?? this.init.perPage ?? 10, this.buildOptions());
	}

	getFullList() {
		return this.pb
			.collection(this.collection)
			.getFullList<QueryResponse<C, Expand>>(this.buildOptions());
	}

	getOne(id: string): Promise<QueryResponse<C, Expand>> {
		return this.pb
			.collection(this.collection)
			.getOne<QueryResponse<C, Expand>>(id, this.buildOptions());
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
	return filter.fields.map((f) => `${f} ~ ${QUOTE}${filter.text}${QUOTE}`).join(' || ');
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

//

type PocketbaseListOptions = Simplify<RecordFullListOptions & RecordListOptions>;

//

export type CollectionExpand<C extends CollectionName> = KeyOf<CollectionExpands[C]>[];

type ResolveCollectionExpand<C extends CollectionName, E extends CollectionExpand<C>> = Partial<
	Pick<CollectionExpands[C], E[number]>
>;

export type QueryResponse<
	C extends CollectionName,
	Expand extends CollectionExpand<C> = never
> = CollectionResponses[C] &
	Simplify<{
		expand?: ResolveCollectionExpand<C, Expand>;
	}>;
