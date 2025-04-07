import { pb } from '@/pocketbase/index.js';
import { z } from 'zod';
import { workflowResponse as SAMPLE_WORKFLOW_RESPONSE } from './components/data';
import { error } from '@sveltejs/kit';
import type { HistoryEvent } from '@forkbombeu/temporal-ui';

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

	//

	const historyResponse = await pb.send(`${basePath}/history`, {
		method: 'GET',
		fetch
	});
	const historyResponseValidation = rawHistoryResponseSchema.safeParse(historyResponse);
	if (!historyResponseValidation.success) {
		error(500, { message: 'Failed to parse workflow response' });
	}

	//

	return {
		workflowId: workflow_id,
		workflowResponse: SAMPLE_WORKFLOW_RESPONSE,
		eventHistory: historyResponseValidation.data as HistoryEvent[]
	};
};

//

const rawWorkflowResponseSchema = z.record(z.unknown());

const rawHistoryResponseSchema = z.array(z.record(z.unknown()));
