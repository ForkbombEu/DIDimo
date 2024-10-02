// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import PocketBase from 'pocketbase';
import fs from 'fs';
import 'dotenv/config';

const { PB_ADMIN_USER, PB_ADMIN_PASS, PUBLIC_POCKETBASE_URL } = process.env;
const pb = new PocketBase(PUBLIC_POCKETBASE_URL);

await pb.admins.authWithPassword(PB_ADMIN_USER, PB_ADMIN_PASS);

// after the above you can also access the auth data from the authStore
const getSchema = async () => {
	return JSON.stringify(await pb.collections.getFullList(), null, 4);
};

const filePath = './src/lib/pocketbase/schema/db_schema.json';

fs.writeFileSync(filePath, await getSchema());
console.log(`📦 Created pocketbase definitions at ${filePath}`);
