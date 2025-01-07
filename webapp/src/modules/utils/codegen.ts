import type { TypedPocketBase } from '@/pocketbase/types';
import PocketBase from 'pocketbase';
import 'dotenv/config';
import assert from 'node:assert';
import prettier from 'prettier';

//

export const GENERATED = 'generated';
export const EXPORT_TYPE = 'export type ';
export const SEPARATOR = '/* ------------------ */';

//

export async function initAdminPocketbase() {
	const { PB_ADMIN_USER, PB_ADMIN_PASS, PUBLIC_POCKETBASE_URL } = process.env;
	assert(PB_ADMIN_USER);
	assert(PB_ADMIN_PASS);
	assert(PUBLIC_POCKETBASE_URL);

	const pb = new PocketBase(PUBLIC_POCKETBASE_URL) as TypedPocketBase;
	await pb.collection('_superusers').authWithPassword(PB_ADMIN_USER, PB_ADMIN_PASS);
	return pb;
}

export async function formatCode(
	code: string,
	options: prettier.Options = { parser: 'typescript' }
) {
	const formatOptions = await prettier.resolveConfig(import.meta.url, { editorconfig: true });
	const formattedCode = await prettier.format(code, {
		...formatOptions,
		...options
	});
	return formattedCode;
}

export function logCodegenResult(subject: string, filePath: string) {
	console.log('');
	console.log(`📦 Generated ${subject} in: ${filePath}`);
}
