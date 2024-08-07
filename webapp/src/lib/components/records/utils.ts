import { getCollectionSchema } from '$lib/pocketbase/schema';
import { Collections } from '$lib/pocketbase/types';
import type { PBResponse } from '$lib/utils/types';

//

export function excludeIdsFilter(ids: string[]) {
	return ids.map((id) => `id != '${id}'`).join(' && ');
}

export function searchTextFilter(collection: string, text: string) {
	return getCollectionFieldNames(collection)
		.map((f) => `${f} ~ "${text}"`)
		.join(' || ');
}

//

export function mergeFilters(...filters: Array<string | undefined>): string | undefined {
	const validFilters: Array<string> = filters.filter(isNonEmptyString);
	if (validFilters.length == 1) return validFilters[0];
	else if (validFilters.length > 1) return validFilters.map((f) => `(${f})`).join(' && ');
	else return undefined;
}

function isNonEmptyString(s: string | undefined): s is string {
	return typeof s == 'string' && Boolean(s);
}

//

function getCollectionFieldNames(collection: string | Collections): string[] {
	const fieldNames: string[] = [];

	if (collection == '_pb_users_auth_' || collection == Collections.Users) {
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

//

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
