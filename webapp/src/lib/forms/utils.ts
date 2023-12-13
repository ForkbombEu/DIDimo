import { z } from 'zod';

export type ZodFileOptions = {
	types?: string[];
	size?: number;
};

export function zodFile(options: ZodFileOptions = {}) {
	const { size, types } = options;

	let schema = z.custom((v) => v instanceof File, `Input is not a file`);

	if (size)
		schema = schema.refine((v) => (v as File).size < size, `File size exceeds ${size} bytes`);
	if (types)
		schema = schema.refine(
			(v) => types.includes((v as File).type),
			`File type not: ${types.join(', ')}`
		);

	return schema;
}
