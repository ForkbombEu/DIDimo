import { pb } from '@/pocketbase';

export const load = async () => {
	const workflows = await pb.send('/api/workflows', {
		method: 'GET'
	});

	console.log(workflows);

	return {
		workflows
	};
};
