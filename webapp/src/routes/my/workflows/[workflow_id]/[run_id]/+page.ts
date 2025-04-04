import { pb } from '@/pocketbase/index.js';
import {
	workflowResponse as SAMPLE_WORKFLOW_RESPONSE,
	eventHistory as SAMPLE_EVENT_HISTORY
} from './components/data';

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
		workflowResponse: SAMPLE_WORKFLOW_RESPONSE,
		eventHistory: SAMPLE_EVENT_HISTORY
	};
};
