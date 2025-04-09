// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import type { RecordFullListOptions, RecordListOptions } from 'pocketbase';
import type { Simplify } from 'type-fest';

export type PocketbaseListOptions = Simplify<RecordFullListOptions & RecordListOptions>;
