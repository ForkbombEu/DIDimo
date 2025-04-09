// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch }) => {
	const record = (
		await pb.collection('z_test_collection').getFullList({ fetch, requestKey: null })
	).at(0);
	if (!record) error(404);
	return { record };
};
