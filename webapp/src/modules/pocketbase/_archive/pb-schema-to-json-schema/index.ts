// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// import { Effect, pipe, Record } from 'effect';
// import fs from 'node:fs';
// import path from 'node:path';

// import { convertDbConfigToSchemas } from './procedures';
// import type { Plugin } from 'vite';
// import { JSONSchema } from '@effect/schema';
// import { getCurrentWorkingDirectory } from '@/utils/fs';

// import prettier from 'prettier';

// //

// export function convertPbSchemaToJsonSchemaPlugin(): Plugin {
// 	return {
// 		name: 'save_database_index',
// 		buildStart: () => {
// 			convertPbSchemaToJsonSchema();
// 		},
// 		handleHotUpdate: () => {
// 			convertPbSchemaToJsonSchema();
// 		}
// 	};
// }

// convertPbSchemaToJsonSchema();

// function convertPbSchemaToJsonSchema() {
// 	pipe(
// 		convertDbConfigToSchemas(),
// 		Effect.map((schemas) =>
// 			pipe(
// 				schemas,
// 				Record.fromEntries,
// 				Record.map(JSONSchema.make),
// 				(schema) => JSON.stringify(schema, null, 4),
// 				(schema) => `export default ${schema} as const`,
// 				(schema) => prettier.format(schema, { parser: 'babel-ts' }),
// 				(formatPromise) => Effect.promise(() => formatPromise),
// 				Effect.map((databaseSchema) => {
// 					const schemaPath = path.join(
// 						getCurrentWorkingDirectory(import.meta.url),
// 						`schema.generated.ts`
// 					);
// 					fs.writeFileSync(schemaPath, databaseSchema);
// 				})
// 			)
// 		),
// 		Effect.runPromise
// 	);
// }
