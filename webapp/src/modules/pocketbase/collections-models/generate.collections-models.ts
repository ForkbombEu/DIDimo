import { type CollectionModel } from 'pocketbase';
import fs from 'fs';
import 'dotenv/config';
import path from 'node:path';
import { EXPORT_TYPE, formatCode, GENERATED, logCodegenResult, SEPARATOR } from '@/utils/codegen';
import { pipe, Array as A, Record } from 'effect';
import { capitalize, merge } from 'lodash';
import JsonToTS from 'json-to-ts';
import pb_schema from '../../../../../migrations/pb_schema.json';

/* Setup */

const models = pb_schema as CollectionModel[];

/* Codegen */

const COLLECTION_FIELD = `CollectionField`;
const COLLECTION_MODEL = `CollectionModel`;
const IMPORT_STATEMENTS = `
import type { ${COLLECTION_MODEL} } from 'pocketbase'
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

const ANY_COLLECTION_FIELD = `Any${COLLECTION_FIELD}`;
const anySchemaField = `export type ${ANY_COLLECTION_FIELD} = ${schemaFieldOptionsTypesData.map(({ name }) => name).join(' | ')}`;

const ANY_COLLECTION_MODEL = `AnyCollectionModel`;
const anyCollectionModel = `export type ${ANY_COLLECTION_MODEL} = Simplify<SetFieldType<${COLLECTION_MODEL}, 'schema', ${ANY_COLLECTION_FIELD}[]>>;`;

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
		models.flatMap((model) => model.fields),
		A.map((field) => field.type),
		A.dedupe
	);
}

function createFieldOptionsTypeData(
	fieldType: string,
	models: CollectionModel[]
): GeneratedTypeData {
	const typeName = capitalize(fieldType) + COLLECTION_FIELD;
	return pipe(
		models.flatMap((m) => m.fields).filter((f) => f.type == fieldType),
		// merging data in a single object
		// somehow `required` is not present on some sytem fields, we add it here
		(fieldsSchemas) => merge({ required: false }, ...fieldsSchemas),
		// converting to ts
		(data) => JsonToTS(data, { useTypeAlias: true, rootName: typeName })[0],
		//
		(code) => {
			const newCode = code
				.replace('type: string;', `type: "${fieldType}";`)
				.replace('any[]', 'string[]');
			return `export ${newCode}`;
		},
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

//

export function pipeLog<T>(data: T): T {
	console.log(data);
	return data;
}
