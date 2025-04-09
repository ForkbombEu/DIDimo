// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { pb } from '@/pocketbase/index.js';

export const load = async ({ params }) => {
	const issuer = await pb.collection('credential_issuers').getOne(params.issuer_id);
	return { issuer };
};
