// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { pb } from '@/pocketbase';
import { error } from '@sveltejs/kit';
import { z } from 'zod';

//

export const load = async ({ fetch }) => {
	const data = await pb.send('/api/workflows', {
		method: 'GET',
		fetch
	});

	const workflows = genericResponseSchema.safeParse(data);
	if (!workflows.success) {
		error(500, {
			message: 'Failed to parse response'
		});
	}

	return {
		executions: workflows.data.executions
	};
};

const workflowExecutionSchema = z.object({
	execution: z.object({
		runId: z.string(),
		workflowId: z.string()
	}),
	executionTime: z.string(),
	memo: z.record(z.unknown()),
	rootExecution: z.object({
		runId: z.string(),
		workflowId: z.string()
	}),
	startTime: z.string(),
	endTime: z.string().optional(),
	status: z.string(),
	taskQueue: z.string(),
	type: z.object({
		name: z.string()
	})
});

const genericResponseSchema = z.object({
	executions: z.array(workflowExecutionSchema)
});
