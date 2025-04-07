// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// import type { Schema as S } from '@effect/schema'; // eslint-disable-line @typescript-eslint/no-unused-vars
// import type { FieldType, AnyCollectionModel } from '@/pocketbase/collections-models/types';

// /* Filters */

// // -- effect

// export type SchemaFilter<S extends S.Schema.Any> = (self: S) => S.filter<S>;

// export type SchemaFilterCreator<S extends S.Schema.Any, V> = (value: V) => (self: S) => S.filter<S>;

// // -- pocketbase

// export type FieldSchemaFiltersConfig<S extends S.Schema.Any> = Record<
// 	string,
// 	SchemaFilterCreator<S, unknown>
// >;

// export type ArrayFieldSchemaFiltersConfig = Record<
// 	string,
// 	SchemaFilterCreator<S.Schema<ReadonlyArray<unknown>, ReadonlyArray<unknown>, unknown>, unknown>
// >;

// export type FieldSchemaConfig<S extends S.Schema.Any> = {
// 	schema: S;
// 	filters?: FieldSchemaFiltersConfig<S>;
// };

// export type FieldTypeToSchemaConfig = Record<FieldType, FieldSchemaConfig<S.Schema.Any>>;

// /* Effect */

// export type ConverterConfig = {
// 	readonly fieldTypeToBaseSchema: FieldTypeToSchemaConfig;
// 	readonly arrayFieldSchemaFilters: ArrayFieldSchemaFiltersConfig;
// 	readonly pocketbaseConfig: AnyCollectionModel[];
// };
