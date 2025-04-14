// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { fileURLToPath } from 'node:url';
import path from 'node:path';

export function getCurrentWorkingDirectory(fileUrl: string) {
	const __filename = fileURLToPath(fileUrl);
	const __dirname = path.dirname(__filename);
	return __dirname;
}
