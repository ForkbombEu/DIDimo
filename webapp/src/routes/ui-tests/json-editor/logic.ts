import { z, type ZodTypeAny, type ZodRawShape } from 'zod';
import { getExceptionMessage } from '@/utils/errors';
import { Record as R } from 'effect';
import { pb } from '@/pocketbase';

//

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

export type FieldConfig = {
	field_name: string;
	credimi_id: string;
	i18_label: string;
	i18n_description: string;
	field_type: 'string' | 'object';
};

export type TestInput = {
	format: 'json' | 'variables';
	data: Record<string, unknown>;
};

//

export const jsonObjectStringSchema = z.string().superRefine((v, ctx) => {
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
			if (f.field_type == 'string') {
				schema = z.string().nonempty();
			} else if (f.field_type == 'object') {
				schema = jsonObjectStringSchema;
			} else {
				throw new Error(`Invalid field type: ${f.field_type}`);
			}
			return [f.credimi_id, schema];
		})
	);

	return z.object(schemaRawShape);
}

//

export function getVariables(testId: string, filenames: string[]) {
	return pb.send('/api/conformance-checks/configs/placeholders-by-filenames', {
		method: 'POST',
		body: {
			test_id: '', //  TODO - Empty string is mandatory right now
			filenames
		}
	});
}

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
