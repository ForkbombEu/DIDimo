<script lang="ts">
	import {
		workflowRun,
		fullEventHistory,
		currentEventHistory,
		toWorkflowExecution,
		WorkflowHistoryLayout,
		toEventHistory,
		eventViewType
	} from '@temporalio/ui';

	import type { WorkflowRunWithWorkers, WorkflowEvent } from '@temporalio/ui';
	import '@temporalio/ui/dist/_css/temporal.css';

	import { workflowResponse, eventHistory } from './data';
	import { onMount } from 'svelte';

	$workflowRun = { ...$workflowRun, workflow: toWorkflowExecution(workflowResponse) };
	$fullEventHistory = toEventHistory(eventHistory.history.events);
	$currentEventHistory = toEventHistory(eventHistory.history.events);
</script>

<WorkflowHistoryLayout></WorkflowHistoryLayout>
