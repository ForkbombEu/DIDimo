import type { PBResponse, StringKeys } from '$lib/utils/types';

export type RecordInputOptions<R extends PBResponse> = {
	name: string | undefined;
	placeholder: string | undefined;
	displayFields: StringKeys<R>[];
	formatRecord?: (record: R) => string;
	excludeIds: string[];
	required: boolean;
	disabled: boolean;
	filter?: string | undefined;
	expand?: string | undefined;
};
