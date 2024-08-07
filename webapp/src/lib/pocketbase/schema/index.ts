import jsonSchema from './db_schema.json';
import { Collections } from '../types';
import { FieldType, type CollectionSchema, type FieldSchema } from './types';

export * from './types';

/* Main */

export function getCollectionSchema(collection: Collections | string) {
	return jsonSchema.find((e) => e.name === collection || e.id === collection) as
		| CollectionSchema
		| undefined;
}

//

export const schemas = createSchemasRecord();

function createSchemasRecord(): Record<Collections, CollectionSchema> {
	const schemas: Record<string, CollectionSchema> = {};

	for (const collectionName of Object.values(Collections)) {
		const schema = getCollectionSchema(collectionName);
		if (!schema) throw new Error(`Collection schema not found for ${collectionName}`);
		schemas[collectionName] = schema;
	}

	return schemas as Record<Collections, CollectionSchema>;
}

/* Helpers */

export function isArrayField(fieldSchema: FieldSchema): boolean {
	const type = fieldSchema.type;
	if (type !== FieldType.SELECT && type !== FieldType.RELATION && type !== FieldType.FILE)
		return false;
	if (fieldSchema.options.maxSelect === 1) return false;
	else return true;
}
