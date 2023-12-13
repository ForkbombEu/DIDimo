import type { ValueOf } from '$lib/utils/types';

export const RecordActions = {
	CREATE: 'create',
	DELETE: 'delete',
	EDIT: 'edit',
	SELECT: 'select',
	SHARE: 'share'
} as const;

export type RecordAction = ValueOf<typeof RecordActions>;
export const RecordActionsArray = Object.values(RecordActions);
