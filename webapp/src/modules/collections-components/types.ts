import type { CollectionName } from '@/pocketbase/collections-models';
import type { ExpandQueryOption, PocketbaseQueryOptions, QueryResponse } from '@/pocketbase/query';
import type { CollectionRecords } from '@/pocketbase/types';
import type { RecordPresenter } from './utils';
import type { ControlAttrs } from 'formsnap';

//

export type CollectionInputRecordProps<
	C extends CollectionName,
	Expand extends ExpandQueryOption<C> = never
> = {
	queryOptions?: Partial<PocketbaseQueryOptions<C, Expand>>;
	displayFields?: (keyof CollectionRecords[C])[] | undefined;
	displayFn?: RecordPresenter<QueryResponse<C, Expand>> | undefined;
};

export type CollectionInputProps<
	C extends CollectionName,
	Expand extends ExpandQueryOption<C> = never
> = {
	collection: C;
	disabled?: boolean;
	label?: string | undefined;
	placeholder?: string | undefined;
	onSelect?: (record: QueryResponse<C, Expand>) => void;
	clearValueOnSelect?: boolean;
	controlAttrs?: ControlAttrs;
} & CollectionInputRecordProps<C, Expand>;
