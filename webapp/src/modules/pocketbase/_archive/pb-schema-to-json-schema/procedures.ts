// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// import { Option as O, pipe, Record as R, Array as A, Effect, Tuple } from 'effect';
// import { Schema as S } from '@effect/schema';

// import { FieldType as FT, type FieldConfig } from '@/pocketbase/collections-config/types';
// import { isArrayField } from '@/pocketbase/collections-config/utils';

// import { type SchemaFilter, type FieldSchemaFiltersConfig } from './types';
// import { config } from './config';

// // -- Converters

// export function convertDbConfigToSchemas() {
// 	return pipe(
// 		config.pocketbaseConfig,
// 		A.map((collectionConfig) =>
// 			pipe(
// 				convertFieldConfigsToObjectSchema(collectionConfig.schema),
// 				Effect.map((schema) => Tuple.make(collectionConfig.name, schema))
// 			)
// 		),
// 		Effect.all
// 	);
// }

// export function convertFieldConfigsToObjectSchema(fieldsConfig: FieldConfig[]) {
// 	return pipe(
// 		fieldsConfig,
// 		A.map((fieldConfig) =>
// 			pipe(
// 				convertFieldConfigToSchema(fieldConfig),
// 				// Create entries to store in record
// 				Effect.map((schema) => Tuple.make(fieldConfig.name, schema))
// 			)
// 		),
// 		Effect.all,
// 		Effect.map(R.fromEntries),
// 		Effect.map(S.Struct)
// 	);
// }

// export function convertFieldConfigToSchema(fieldConfig: FieldConfig) {
// 	return pipe(
// 		// -- base schema
// 		getBaseSchema(fieldConfig),
// 		// -- base filters
// 		Effect.flatMap((schema) =>
// 			pipe(
// 				getBaseSchemaFilters(fieldConfig),
// 				Effect.map((filters) => applyFiltersToSchema(schema, filters))
// 			)
// 		),
// 		// -- array schema
// 		Effect.flatMap((schema) =>
// 			Effect.if(isArrayField(fieldConfig), {
// 				onFalse: () => Effect.succeed(schema),
// 				onTrue: () =>
// 					pipe(S.Array(schema), (arraySchema) =>
// 						pipe(
// 							getArraySchemaFilters(fieldConfig),
// 							(arrayFilters) => applyFiltersToSchema(arraySchema, arrayFilters),
// 							Effect.succeed
// 						)
// 					)
// 			})
// 		),
// 		// -- required
// 		Effect.flatMap((schema) =>
// 			Effect.if(fieldConfig.required, {
// 				onTrue: () => Effect.succeed(schema),
// 				onFalse: () => Effect.succeed(S.optional(schema))
// 			})
// 		)
// 	);
// }

// // -- Getters

// function getBaseSchema(fieldConfig: FieldConfig) {
// 	return pipe(
// 		getBaseSchemaConfig(fieldConfig),
// 		Effect.map((baseSchemaConfig) => baseSchemaConfig.schema)
// 	);
// }

// function getBaseSchemaFilters(fieldConfig: FieldConfig) {
// 	return pipe(
// 		getBaseSchemaConfig(fieldConfig),
// 		Effect.map((baseSchemaConfig) =>
// 			pipe(
// 				baseSchemaConfig.filters,
// 				O.fromNullable,
// 				O.getOrElse(() => ({}) as FieldSchemaFiltersConfig<S.Schema.Any>)
// 			)
// 		),
// 		Effect.map((filtersConfig) =>
// 			pipe(
// 				Object.entries(fieldConfig.options),
// 				A.filter(([optionName]) => optionName in filtersConfig),
// 				A.filter(([, optionValue]) => Boolean(optionValue)),
// 				A.map(([optionName, optionValue]) =>
// 					pipe(filtersConfig, R.get(optionName), O.getOrThrow, (filterCreator) =>
// 						filterCreator(optionValue)
// 					)
// 				)
// 			)
// 		)
// 	);
// }

// function getArraySchemaFilters(fieldConfig: FieldConfig) {
// 	return pipe(
// 		Object.entries(fieldConfig.options),
// 		A.filter(([optionName]) => optionName in config.arrayFieldSchemaFilters),
// 		A.map(([optionName, optionValue]) =>
// 			pipe(config.arrayFieldSchemaFilters, R.get(optionName), O.getOrThrow, (filterCreator) =>
// 				filterCreator(optionValue)
// 			)
// 		)
// 	);
// }

// function getBaseSchemaConfig(fieldConfig: FieldConfig) {
// 	return R.get(config.fieldTypeToBaseSchema, fieldConfig.type as FT);
// }

// // -- Utils

// function applyFiltersToSchema(schema: S.Schema.Any, filters: SchemaFilter<S.Schema.Any>[]) {
// 	return A.reduce(filters, schema, (schema, filter) => filter(schema));
// }
