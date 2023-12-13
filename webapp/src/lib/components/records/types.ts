import type { PBResponse, StringKeys } from '$lib/utils/types';

export type RecordInputOptions<R extends PBResponse> = {
	name: string | undefined;
	placeholder: string | undefined;
	displayFields: StringKeys<R>[];
	excludeIds: string[];
	required: boolean;
	disabled: boolean;
};
