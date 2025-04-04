import { pb } from '@/pocketbase/index.js';

export const load = async ({ params }) => {
	const { workflow_id, run_id } = params;

	const response = await pb.send(`/api/workflows/${workflow_id}/${run_id}`, {
		method: 'GET'
	});

	console.log(response);

	return {
		workflow_id
	};
};
