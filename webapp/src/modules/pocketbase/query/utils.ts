import type { RecordFullListOptions, RecordListOptions } from 'pocketbase';
import type { Simplify } from 'type-fest';

export type PocketbaseListOptions = Simplify<RecordFullListOptions & RecordListOptions>;
