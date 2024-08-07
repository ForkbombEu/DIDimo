import { isFile, isFileArray } from '$lib/forms/fields';
import { isArrayField } from '$lib/pocketbase/schema';
import {
	FieldType,
	type FieldSchema,
	type CollectionSchema,
	type FieldOptions
} from '$lib/pocketbase/schema/types';

//

type InitialData = Record<string, unknown>;

/**
 * File handling
 *
 * File field zod schema expects a file
 * but InitialData coming from PocketBase is a string
 *
 * We solve it this way:
 * - 1. Store the original filenames
 * - 2. Convert the strings to "placeholder" files
 * - 3. When submitting the form, match the new files with the original filenames
 */

/** */

// 1. Store original data

export function getFileFieldInitialData(field: FieldSchema, initialData: InitialData) {
	if (field.type == FieldType.FILE && field.name in initialData) {
		const data = initialData[field.name] as string | string[];
		if (!isArrayField(field) && typeof data == 'string') {
			if (data) return [data];
			else return [];
		} else if (isArrayField(field) && Array.isArray(data)) {
			return data;
		}
	}
}

export function getFileFieldsInitialData(collection: CollectionSchema, initialData: InitialData) {
	const fileFieldsInitialData: Record<string, string[]> = {};
	for (const field of collection.schema) {
		const data = getFileFieldInitialData(field, initialData);
		if (data) fileFieldsInitialData[field.name] = data;
	}
	return fileFieldsInitialData;
}

// 2. Convert strings to "placeholder" files

function mockFile(filename: string, fieldOptions: FieldOptions = {}) {
	let fileOptions = undefined;
	const mimeTypes = fieldOptions.mimeTypes;
	if (Array.isArray(mimeTypes) && mimeTypes.length > 0) {
		fileOptions = { type: mimeTypes[0] };
	}
	const mockFile = new File([], filename, fileOptions);
	return mockFile;
}

export function mockFileFieldInitialData(field: FieldSchema, initialData: InitialData) {
	if (field.type == FieldType.FILE && field.name in initialData) {
		const data = initialData[field.name] as string | string[];
		if (!isArrayField(field) && typeof data == 'string') {
			if (data) return mockFile(data, field.options);
			else return undefined;
		} else if (isArrayField(field) && Array.isArray(data)) {
			return data.map((item) => mockFile(item, field.options));
		}
	}
}

export function mockFileFieldsInitialData(collection: CollectionSchema, initialData: InitialData) {
	const data = { ...initialData };
	for (const field of collection.schema) {
		if (field.type == FieldType.FILE && field.name in initialData) {
			data[field.name] = mockFileFieldInitialData(field, initialData);
		}
	}
	return data;
}

// 3. Match the new files with the original filenames

export function cleanFormDataFiles(
	formData: Record<string, unknown>,
	fileFieldsInitialData: Record<string, string[]>
) {
	const data = { ...formData };

	for (const [field, initialFilenames] of Object.entries(fileFieldsInitialData)) {
		const fieldValue = data[field];

		if (fieldValue === undefined || fieldValue === null) {
			continue;
		}
		//
		else if (isFile(fieldValue)) {
			const isFileOld = initialFilenames.includes(fieldValue.name);
			if (isFileOld) delete data[field];
		}
		//
		else if (isFileArray(fieldValue)) {
			const allFilenames = fieldValue.map((file) => file.name);
			const newFiles = fieldValue.filter((file) => !initialFilenames.includes(file.name));
			const filesToRemove = initialFilenames.filter((filename) => !allFilenames.includes(filename));

			if (newFiles.length === 0) delete data[field];
			else data[field] = newFiles;

			if (filesToRemove.length > 0) data[`${field}-`] = filesToRemove;
		}
	}

	return data;
}
