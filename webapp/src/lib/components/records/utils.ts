import { getCollectionSchema } from '$lib/pocketbase/schema';
import { Collections } from '$lib/pocketbase/types';
import type { PBResponse } from '$lib/utils/types';

export function filterStringArray(
	field: string,
	operator: string,
	joinOperator: string,
	values: string[]
) {
	if (values.length === 0) return '';
	const filterString = values.map((v) => `${field} ${operator} '${v}'`).join(` ${joinOperator} `);
	return `(${filterString})`;
}

export function excludeStringArray(field: string, values: string[]) {
	return filterStringArray(field, '!=', '&&', values);
}

export function getCollectionFieldNames(collection: string | Collections): string[] {
	const fieldNames: string[] = [];

	if (collection == '_pb_users_auth_' || collection == Collections.Users) {
		fieldNames.push('username');
		fieldNames.push('email');
	}

	const collectionSchema = getCollectionSchema(collection);
	if (collectionSchema) {
		collectionSchema.schema.forEach((field) => {
			fieldNames.push(field.name);
		});
	}

	return fieldNames;
}

export function createRecordLabel<R extends PBResponse>(
	record: R,
	labelFields: string[] = []
): string {
	const fields: string[] = [];

	if (labelFields.length > 0) {
		fields.push(...labelFields);
	} else {
		fields.push(...getCollectionFieldNames(record.collectionName));
	}

	return fields
		.map((f) => record[f])
		.filter((f) => Boolean(f))
		.join(' | ');
}
