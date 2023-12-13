import type { PBResponse } from '$lib/utils/types';
import type { FieldComponent } from './fieldComponents/fieldComponentRenderer.svelte';
import type { RecordAction } from './recordActions';

export type FieldsComponents<RecordGeneric extends PBResponse = PBResponse> = Record<
	string,
	FieldComponent<RecordGeneric>
>;

export type ViewAction = Exclude<RecordAction, 'create'>;
