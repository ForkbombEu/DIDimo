import { pipe } from 'effect';
import z from 'zod';
import type { SchemaFields } from '@/pocketbase/collections-models';
import { getJsonDataSize } from '@/utils/other';
import { isBefore, isAfter, isValid, parseISO } from 'date-fns';
import { m } from '@/i18n';
import { zodFileSchema } from '@/utils/files';

/* Field Config -> Zod Type */

type SchemaFieldToZodTypeMap = {
	[Type in keyof SchemaFields]: (fieldSchema: SchemaFields[Type]) => z.ZodTypeAny;
};

export const schemaFieldToZodTypeMap: SchemaFieldToZodTypeMap = {
	text: (config) => {
		const { max, min, pattern } = config.options;
		let s = z.string();
		if (max) s = s.max(max);
		if (min) s = s.min(min);
		if (pattern) {
			// // Add a "|" pipe to the regex to allow for empty string (Ciscoheat suggestion)
			// const maybeOptionalPattern = config.required ? pattern : `|${pattern}`;
			// TODO - Check if it is needed still
			s = s.regex(new RegExp(pattern), m.Value_does_not_match_regex_pattern({ pattern }));
		}
		return s;
	},

	bool: () => {
		return z.boolean();
	},

	email: ({ options }) => {
		const { exceptDomains, onlyDomains } = options;
		return pipe(z.string().email(), (zodEmail) =>
			validateDomains(zodEmail, exceptDomains, onlyDomains)
		);
	},

	file: ({ options }) => {
		const { mimeTypes, maxSize } = options;
		const mimes = mimeTypes as string[] | undefined;
		return zodFileSchema({ mimeTypes: mimes, maxSize });
	},

	date: ({ options }) => {
		const { min, max } = options;
		return z
			.string()
			.refine(
				(string) => isValid(parseISO(string)),
				(value) => ({ message: `${value} is not a ISO date string` })
			)
			.refine(
				(date) => (min ? isAfter(date, min) : true),
				(value) => ({ message: `${value} is before ${min}` })
			)
			.refine(
				(date) => (max ? isBefore(date, max) : true),
				(value) => ({ message: `${value} is after ${max}` })
			);
	},

	json: ({ options }) => {
		const { maxSize } = options;
		return z.unknown().refine((json) => {
			if (maxSize) return getJsonDataSize(json) < maxSize;
			else return true;
		}, `Json size is bigger than ${maxSize} bytes`);
	},

	relation: () => {
		return z.string();
	},

	number: ({ options }) => {
		const { min, max, noDecimal } = options;
		let s = z.number();
		if (min) s = s.min(min);
		if (max) s = s.max(max);
		if (noDecimal) s = s.int();
		return s;
	},

	select: ({ options }) => {
		const { values } = options;
		if (!values) throw new SelectSchemaFieldNoOptionsError();
		return z.string().refine((s) => values.includes(s));
	},

	editor: () => {
		return z.string();
	},

	url: ({ options }) => {
		const { exceptDomains, onlyDomains } = options;
		return pipe(z.string().url(), (zodUrl) => validateDomains(zodUrl, exceptDomains, onlyDomains));
	}
};

class SelectSchemaFieldNoOptionsError extends Error {}

//

function validateDomains(
	zodString: z.ZodString,
	exceptDomains: readonly string[] | undefined = undefined,
	onlyDomains: readonly string[] | undefined = undefined
) {
	let s: z.ZodString | z.ZodEffects<z.ZodString> | z.ZodEffects<z.ZodEffects<z.ZodString>> =
		zodString;

	if (onlyDomains?.length) {
		s = s.refine(
			(string) => onlyDomains.some((domain) => string.includes(domain)),
			m.URL_is_not_in_allowed_domains_list() + ': ' + (onlyDomains ?? []).join(', ')
		);
	}

	if (exceptDomains?.length) {
		s = s.refine(
			(string) => exceptDomains.every((domain) => !string.includes(domain)),
			m.URL_is_in_forbidden_domains_list() + ': ' + (exceptDomains ?? []).join(', ')
		);
	}

	return s;
}
