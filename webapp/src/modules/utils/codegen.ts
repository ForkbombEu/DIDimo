import 'dotenv/config';
import prettier from 'prettier';

import sqlite3 from 'sqlite3';
import { open } from 'sqlite';

//

export function openDb() {
	return open({
		filename: process.env.DATA_DB_PATH || "../pb_data/data.db",
		driver: sqlite3.Database
	});
}

//

export const GENERATED = 'generated';
export const EXPORT_TYPE = 'export type ';
export const SEPARATOR = '/* ------------------ */';

//

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
