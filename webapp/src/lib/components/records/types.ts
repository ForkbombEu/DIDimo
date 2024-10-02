// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

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
