import type { BaseSystemFields } from '$lib/pocketbase/types';
import type { XCircle } from 'svelte-heros-v2';

/* TS utils */

export type ValueOf<T> = T[keyof T];

export type ArrayExtract<T> = T extends (infer U)[] ? U : T;
// Util to extract type from T | T[]

export type StringKeys<T> = Extract<keyof T, string>;

/* Components */

export type Link = {
	href: string;
	text: string;
};

export type IconComponent = typeof XCircle;

/* Pocketbase generics */

export type PBRecord = Record<string, unknown>;
export type PBResponse = PBRecord & BaseSystemFields<any>;

export type ExtractPBRecord<T> = T extends infer R & BaseSystemFields<unknown> ? R : never;
export type ExtractPBExpand<T> = T extends unknown & BaseSystemFields<infer E> ? E : never;
