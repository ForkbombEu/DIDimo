import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';
import { z } from 'zod';

export const load = async ({ fetch }) => {
	const data = await pb.send('/api/workflows', {
		method: 'GET',
		fetch
	});

	const workflows = WorkflowResponseSchema.safeParse(data);
	if (!workflows.success) {
		error(500, {
			message: 'Failed to parse workflows'
		});
	}

	return {
		executions: workflows.data.executions
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
