import type { CollectionFormOptions } from '@/collections-components/form';
import type { CollectionName } from '@/pocketbase/collections-models';
import type { PocketbaseQueryExpandOption } from '@/pocketbase/query';
import { CollectionManager } from './collectionManager.svelte.js';
import { setupDerivedContext } from '@/utils/svelte-context';
import type { FilterMode } from '@/pocketbase/query/query.js';

//

export type Filter = {
	name: string;
	expression: string;
};

export type FilterGroup = {
	name?: string;
	id: string;
	mode: FilterMode;
	filters: Filter[];
};

export type FiltersOption = FilterGroup | FilterGroup[];

//

export type CollectionManagerContext<
	C extends CollectionName = never,
	Expand extends PocketbaseQueryExpandOption<C> = never
> = {
	manager: CollectionManager<C, Expand>;
	filters: FiltersOption;
	formsOptions: Record<FormPropType, CollectionFormOptions<C>>;
};

type FormPropType = 'base' | 'create' | 'edit';

export const {
	getDerivedContext: getCollectionManagerContext,
	setDerivedContext: setCollectionManagerContext
} = setupDerivedContext<CollectionManagerContext>('cmc');
