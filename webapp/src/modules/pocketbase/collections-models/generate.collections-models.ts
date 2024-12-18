import { type CollectionModel } from 'pocketbase';
import fs from 'fs';
import 'dotenv/config';
import path from 'node:path';
import {
	EXPORT_TYPE,
	formatCode,
	GENERATED,
	initAdminPocketbase,
	logCodegenResult,
	SEPARATOR
} from '@/utils/codegen';
import { pipe, Array as A, String, Record } from 'effect';
import { capitalize, merge } from 'lodash';
import JsonToTS from 'json-to-ts';

/* Setup */

const pb = await initAdminPocketbase();
const models = await pb.collections.getFullList();

/* Codegen */

const SCHEMA_FIELD = `SchemaField`;
const COLLECTION_MODEL = `CollectionModel`;
const IMPORT_STATEMENTS = `
import type { ${SCHEMA_FIELD}, ${COLLECTION_MODEL} } from 'pocketbase'
import type { SetFieldType, Simplify } from 'type-fest';
`;

const schemaFieldTypes = getFieldTypeNames(models);
const schemaFieldType = `${EXPORT_TYPE} SchemaFieldType = ${schemaFieldTypes.map((t) => JSON.stringify(t)).join(' | ')}`;
const schemaFieldOptionsTypesData = schemaFieldTypes.map((f) =>
	createFieldOptionsTypeData(f, models)
);
const schemaFields = `export type SchemaFields = {
    ${schemaFieldOptionsTypesData.map(({ name, key }) => `${key}: ${name}`).join('\n')}
}`;
const anySchemaField = `export type AnySchemaField = ${schemaFieldOptionsTypesData.map(({ name }) => name).join(' | ')}`;

const ANY_COLLECTION_MODEL = `AnyCollectionModel`;
const anyCollectionModel = `export type ${ANY_COLLECTION_MODEL} = Simplify<SetFieldType<${COLLECTION_MODEL}, 'schema', AnySchemaField[]>>;`;

const code = [
	IMPORT_STATEMENTS,
	SEPARATOR,
	schemaFieldType,
	anySchemaField,
	schemaFields,
	...schemaFieldOptionsTypesData.map((data) => data.code),
	SEPARATOR,
	anyCollectionModel,
	collectionName(models),
	sanitizeCollectionsModels(models)
].join('\n\n');

/* Export */

const formattedCode = await formatCode(code);
const filePath = path.resolve(import.meta.dirname, `collections-models.${GENERATED}.ts`);
fs.writeFileSync(filePath, formattedCode);
logCodegenResult('collections models and helper types', filePath);

/* Helper functions */

function sanitizeCollectionsModels(models: CollectionModel[]) {
	// Hiding API rules to reduce leaked information
	const sanitizedModels = models.map(
		Record.map((v, k) => {
			if (k.includes('Rule')) return '';
			else return v;
		})
	);

	return `export const CollectionsModels = ${JSON.stringify(sanitizedModels, null, 2)} as ${ANY_COLLECTION_MODEL}[]`;
}

function collectionName(models: CollectionModel[]): string {
	const names = models.map((m) => m.name);
	return `export type CollectionName = ${names.map((n) => JSON.stringify(n)).join(' | ')}`;
}

//

function getFieldTypeNames(models: CollectionModel[]) {
	return pipe(
		models.flatMap((model) => model.schema),
		A.map((field) => field.type),
		A.dedupe
	);
}

function createFieldOptionsTypeData(
	fieldType: string,
	models: CollectionModel[]
): GeneratedTypeData {
	const typeName = capitalize(fieldType) + SCHEMA_FIELD;
	return pipe(
		models
			.flatMap((m) => m.schema)
			.filter((f) => f.type == fieldType)
			.map((f) => f.options),
		// merging data in a single object
		(fieldsSchemas) => merge({}, ...fieldsSchemas),
		// converting to ts
		(data) => JsonToTS(data, { useTypeAlias: true })[0],
		//
		(code) => {
			const typeOnly = code.split('=')[1].trim();
			const removeEmptyObject = typeOnly == `{\n}` ? 'Record<string,never>' : typeOnly;
			const newCode = `export type ${typeName} = ${SCHEMA_FIELD} & { type: "${fieldType}"; options: Partial<${removeEmptyObject}>; unique: boolean }`;
			return newCode;
		},
		// small fix
		String.replace('any[]', 'string[]'),
		//
		(code) => ({
			code,
			name: typeName,
			key: fieldType
		})
	);
}

type GeneratedTypeData = {
	code: string;
	name: string;
	key: string;
};
