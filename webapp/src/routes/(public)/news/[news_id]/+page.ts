// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { PocketbaseQueryAgent } from '@/pocketbase/query/agent.js';

export const load = async ({ params }) => {
	const news = await new PocketbaseQueryAgent({
		collection: 'news'
	}).getOne(params.news_id);

	return {
		news
	};
};
