import { fileURLToPath } from 'node:url';
import path from 'node:path';

export function getCurrentWorkingDirectory(fileUrl: string) {
	const __filename = fileURLToPath(fileUrl);
	const __dirname = path.dirname(__filename);
	return __dirname;
}
