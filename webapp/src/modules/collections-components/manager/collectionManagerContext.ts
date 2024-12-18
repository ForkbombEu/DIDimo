import type { CollectionFormOptions } from '@/collections-components/form';
import type { CollectionName } from '@/pocketbase/collections-models';
import type { ExpandQueryOption } from '@/pocketbase/query';
import { CollectionManager } from './collectionManager.svelte.js';
import { setupDerivedContext } from '@/utils/svelte-context';

//

export type CollectionManagerContext<
	C extends CollectionName = never,
	Expand extends ExpandQueryOption<C> = never
> = {
	manager: CollectionManager<C, Expand>;
	formsOptions: Record<FormPropType, CollectionFormOptions<C>>;
};

type FormPropType = 'base' | 'create' | 'edit';

export const {
	getDerivedContext: getCollectionManagerContext,
	setDerivedContext: setCollectionManagerContext
} = setupDerivedContext<CollectionManagerContext>('cmc');
