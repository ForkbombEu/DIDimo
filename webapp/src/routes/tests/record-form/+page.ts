import { pb } from '$lib/pocketbase';
import {
	Collections,
	type CrudExampleResponse,
	type FeaturesResponse,
	type UsersResponse
} from '$lib/pocketbase/types';

export const load = async () => {
	try {
		const item = await pb
			.collection(Collections.CrudExample)
			.getOne<CrudExampleResponse>('53q5yve2r8scku7');
		return { item };
	} catch (_) {
		return {};
	}
};
