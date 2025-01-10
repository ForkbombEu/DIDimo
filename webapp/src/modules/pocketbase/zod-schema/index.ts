import z from 'zod';
import { pipe } from 'effect';
import { schemaFieldToZodTypeMap } from './config';
import {
	getCollectionModel,
	isArrayField,
	type AnyCollectionField,
	type CollectionName
} from '@/pocketbase/collections-models';
import { systemFields, type CollectionZodRawShapes } from '../types';

//

export type CollectionZodSchema<C extends CollectionName> = z.ZodObject<CollectionZodRawShapes[C]>;

export function getCollectionFields(collection: CollectionName) {
	const { fields } = getCollectionModel(collection);
	return (fields as AnyCollectionField[]).filter(
		(f) => !(systemFields as unknown as string[]).includes(f.name)
	);
}

export function createCollectionZodSchema<C extends CollectionName>(
	collection: C
): CollectionZodSchema<C> {
	const collectionFields = getCollectionFields(collection);

	const entries = collectionFields.map((fieldConfig) => {
		const zodTypeConstructor = schemaFieldToZodTypeMap[fieldConfig.type] as (
			c: AnyCollectionField
		) => z.ZodTypeAny;

		const zodType = pipe(
			zodTypeConstructor(fieldConfig),

			// Array type handling
			(zodType) => {
				if (isArrayField(fieldConfig)) {
					let s = z.array(zodType);
					const { minSelect, maxSelect } = z
						.object({
							minSelect: z.number().nullish(),
							maxSelect: z.number().nullish()
						})
						.parse(fieldConfig);
					if (minSelect) s = s.min(minSelect);
					if (maxSelect) s = s.max(maxSelect);
					return s;
				} else {
					return zodType;
				}
			},

			// Optional type handling
			(zodType) => {
				if (fieldConfig.required) {
					if (zodType instanceof z.ZodArray) return zodType.nonempty();
					else return zodType;
				} else {
					// Extra check for url: https://github.com/colinhacks/zod/discussions/1254
					if (fieldConfig.type == 'url') return zodType.or(z.literal('')).optional();
					else return zodType.optional();
				}
			}
		);

		return [fieldConfig.name, zodType];
	});

	const rawObject = Object.fromEntries(entries);
	return z.object(rawObject) as CollectionZodSchema<C>;
}
