import { z, type ZodTypeAny, type ZodRawShape } from 'zod';
import { getExceptionMessage } from '@/utils/errors';
import { Record as R } from 'effect';
import { pb } from '@/pocketbase';

/* Schemas */

const sharedFieldSchema = z.object({
	CredimiID: z.string(),
	DescriptionKey: z.string(),
	LabelKey: z.string(),
	Type: z.literal('string').or(z.literal('object'))
});

export type FieldConfig = z.infer<typeof sharedFieldSchema>;

const specificFieldSchema = sharedFieldSchema.extend({
	FieldName: z.string(),
	Example: z.string().optional()
});

export type SpecificFieldConfig = z.infer<typeof specificFieldSchema>;

//

export const stringifiedObjectSchema = z.string().superRefine((v, ctx) => {
	try {
		z.record(z.string(), z.unknown())
			.refine((value) => R.size(value) > 0)
			.parse(JSON.parse(v));
	} catch (e) {
		const message = getExceptionMessage(e);
		ctx.addIssue({
			code: z.ZodIssueCode.custom,
			message: `Invalid JSON object: ${message}`
		});
	}
});

export function createSchemaFromFieldsConfigs(fields: FieldConfig[]) {
	const schemaRawShape: ZodRawShape = Object.fromEntries(
		fields.map((f) => {
			let schema: ZodTypeAny;
			if (f.Type == 'string') {
				schema = z.string().nonempty();
			} else if (f.Type == 'object') {
				schema = stringifiedObjectSchema;
			} else {
				throw new Error(`Invalid field type: ${f.Type}`);
			}
			return [f.CredimiID, schema];
		})
	);

	return z.object(schemaRawShape);
}

const fieldsResponseSchema = z.object({
	normalized_fields: z.array(sharedFieldSchema),
	specific_fields: z.record(
		z.string(),
		z.object({
			content: stringifiedObjectSchema,
			fields: z.array(specificFieldSchema)
		})
	)
});

export type FieldsResponse = z.infer<typeof fieldsResponseSchema>;

export type TestInput = {
	format: 'json' | 'variables';
	data: Record<string, unknown>;
};

export async function getVariables(testId: string, filenames: string[]) {
	const data = await pb.send('/api/conformance-checks/configs/placeholders-by-filenames', {
		method: 'POST',
		body: {
			test_id: '', //  TODO - Empty string is mandatory right now
			filenames
		}
	});
	return fieldsResponseSchema.parse(data);
}

export type StandardWithTestSuites = {
	id: string;
	label: string;
	description: string;
	testSuites: Array<{
		id: string;
		label: string;
		tests: string[];
	}>;
};

//

// let fieldsStatus: Record<string, boolean> = $state({});

// watch(
// 	() => formState.current,
// 	() => {
// 		Promise.all(
// 			fields.map(async (f) => {
// 				const result = await validate(f.credimi_id, { taint: false, update: false });
// 				const isValid = result === undefined || result.length === 0;
// 				return Tuple.make(f.i18_label, isValid);
// 			})
// 		).then((statuses) => {
// 			fieldsStatus = Object.fromEntries(statuses);
// 		});
// 	}
// );
