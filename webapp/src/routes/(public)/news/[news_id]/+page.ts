import { PocketbaseQueryAgent } from '@/pocketbase/query/agent.js';

export const load = async ({ params }) => {
	const news = await new PocketbaseQueryAgent({
		collection: 'news'
	}).getOne(params.news_id);

	return {
		news
	};
};
