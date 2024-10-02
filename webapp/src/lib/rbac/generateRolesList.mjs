// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import fs from 'fs';
import PocketBase from 'pocketbase';
import 'dotenv/config';

//

const objectName = 'OrgRoles';
const typeName = 'OrgRole';
const filePath = './src/lib/rbac/roles.ts';

// Pocketbase setup

const { PB_ADMIN_USER, PB_ADMIN_PASS, PUBLIC_POCKETBASE_URL } = process.env;

const pb = new PocketBase(PUBLIC_POCKETBASE_URL);
await pb.admins.authWithPassword(PB_ADMIN_USER, PB_ADMIN_PASS);

const roles = await pb.collection('orgRoles').getFullList();

// Building file

let file = '';

file += `export const ${objectName} = {\n`;
for (const role of roles) {
	file += `\t${role.name.toUpperCase()}: '${role.name}',\n`;
}
file += '} as const;\n';

file += '\n';
file += `export type ${typeName} = typeof ${objectName} [keyof typeof ${objectName}];`;

fs.writeFileSync(filePath, file);
console.log(`📦 Generated ${objectName} list in: ${filePath}`);
