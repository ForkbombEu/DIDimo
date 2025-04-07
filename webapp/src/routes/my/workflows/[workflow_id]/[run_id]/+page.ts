import { pb } from '@/pocketbase/index.js';

//

export const load = async ({ params, fetch }) => {
	const { workflow_id, run_id } = params;

	const basePath = `/api/workflows/${workflow_id}/${run_id}`;

	// TODO - Make sure that the result of these requests matches the shape of the data imported from `./components/data.ts`
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	const workflowResponse = await pb.send(basePath, {
		method: 'GET',
		fetch
	});
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	const historyResponse = await pb.send(`${basePath}/history`, {
		method: 'GET',
		fetch
	});

	return {
		workflowId: workflow_id,
		workflowResponse,
		eventHistory: { history: { events: [] } }
	};
};
