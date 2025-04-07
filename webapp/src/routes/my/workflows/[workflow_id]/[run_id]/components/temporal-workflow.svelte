<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import {
		workflowRun,
		fullEventHistory,
		currentEventHistory,
		toWorkflowExecution,
		WorkflowHistoryLayout,
		toEventHistory,
		type HistoryEvent
	} from '@forkbombeu/temporal-ui';

	//

	type Props = {
		workflowResponse: Record<string, unknown>;
		eventHistory: { history: { events: HistoryEvent[] } };
	};

	let { workflowResponse, eventHistory }: Props = $props();

	//

	$workflowRun = { ...$workflowRun, workflow: toWorkflowExecution(workflowResponse) };
	$fullEventHistory = toEventHistory(eventHistory.history.events);
	$currentEventHistory = toEventHistory(eventHistory.history.events);

	//

	beforeNavigate(({ cancel, to }) => {
		const pathname = to?.url.pathname;
		if (pathname?.includes('undefined')) cancel();
	});
</script>

<div class="temporal-ui-workflow space-y-4">
	<WorkflowHistoryLayout></WorkflowHistoryLayout>
</div>

<style lang="postcss">
	:global(div > table > tbody > div.text-right.hidden) {
		display: block;
	}

	:global(button.toggle-button[data-testid='download']) {
		display: none;
	}

	:global(button.toggle-button[data-testid='pause']) {
		display: none;
	}

	:global(.temporal-ui-workflow a) {
		@apply !no-underline hover:!cursor-not-allowed hover:!text-inherit;
	}
</style>
