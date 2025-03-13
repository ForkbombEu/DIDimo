<script lang="ts">
	import './temporal.css';

	import {
		workflowRun,
		fullEventHistory,
		currentEventHistory,
		toWorkflowExecution,
		WorkflowHistoryLayout,
		toEventHistory,
		eventViewType
	} from 'temporal-ui';

	import type { WorkflowRunWithWorkers, WorkflowEvent } from 'temporal-ui';

	import { workflowResponse, eventHistory } from './data';
	import { onMount } from 'svelte';

	$workflowRun = { ...$workflowRun, workflow: toWorkflowExecution(workflowResponse) };
	$fullEventHistory = toEventHistory(eventHistory.history.events);
	$currentEventHistory = toEventHistory(eventHistory.history.events);
</script>

<WorkflowHistoryLayout></WorkflowHistoryLayout>
