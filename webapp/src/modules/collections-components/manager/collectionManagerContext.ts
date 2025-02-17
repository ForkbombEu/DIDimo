import type { CollectionFormOptions } from '@/collections-components/form';
import type { CollectionName } from '@/pocketbase/collections-models';
import type { ExpandQueryOption } from '@/pocketbase/query';
import { CollectionManager } from './collectionManager.svelte.js';
import { setupDerivedContext } from '@/utils/svelte-context';
import { z } from 'zod';

//

export const FilterSchema = z.object({
	id: z.string(),
	name: z.string(),
	expression: z.string()
});

export type Filter = z.infer<typeof FilterSchema>;

export const FilterGroupSchema = z.object({
	name: z.string(),
	filters: z.array(FilterSchema)
});

export type FilterGroup = z.infer<typeof FilterGroupSchema>;

//

export type CollectionManagerContext<
	C extends CollectionName = never,
	Expand extends ExpandQueryOption<C> = never
> = {
	manager: CollectionManager<C, Expand>;
	filters: (Filter | FilterGroup)[];
	formsOptions: Record<FormPropType, CollectionFormOptions<C>>;
};

type FormPropType = 'base' | 'create' | 'edit';

export const {
	getDerivedContext: getCollectionManagerContext,
	setDerivedContext: setCollectionManagerContext
} = setupDerivedContext<CollectionManagerContext>('cmc');
