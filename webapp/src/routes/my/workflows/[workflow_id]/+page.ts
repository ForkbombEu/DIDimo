import { error } from '@sveltejs/kit';

export const load = async () => {
	error(404);
};
