import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';
import { z } from 'zod';
import { toWorkflowExecutions, type ListWorkflowExecutionsResponse } from '@forkbombeu/temporal-ui';

export const load = async () => {
	const data = await pb.send('/api/workflows', {
		method: 'GET'
	});

	const workflows = WorkflowResponseSchema.safeParse(data);
	if (!workflows.success) {
		error(500, {
			message: 'Failed to parse workflows'
		});
	}

	const executions: ListWorkflowExecutionsResponse = workflows.data.executions;
	const d = toWorkflowExecutions(executions);

	return {
		workflows: d
	};
};

//

const WorkflowSchema = z.object({
	execution: z.object({
		workflow_id: z.string(),
		run_id: z.string()
	}),
	type: z.object({
		name: z.string()
	}),
	start_time: z.object({
		seconds: z.number(),
		nanos: z.number()
	}),
	status: z.number(),
	execution_time: z.object({
		seconds: z.number(),
		nanos: z.number()
	}),
	memo: z.record(z.never()),
	task_queue: z.string(),
	root_execution: z.object({
		workflow_id: z.string(),
		run_id: z.string()
	})
});

const WorkflowResponseSchema = z.object({
	executions: z.array(WorkflowSchema)
});
