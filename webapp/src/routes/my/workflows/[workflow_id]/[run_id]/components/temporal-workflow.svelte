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
		// pauseLiveUpdates
	} from '@forkbombeu/temporal-ui';

	//

	type Props = {
		workflowResponse: Record<string, unknown>;
		eventHistory: HistoryEvent[];
	};

	let { workflowResponse, eventHistory }: Props = $props();

	//

	const workflow = toWorkflowExecution(workflowResponse);
	/* HACK */
	// canBeTerminated a property of workflow object is a getter that requires a svelte `store` to work
	// by removing it, we can avoid the store dependency and solve a svelte error about state not updating
	Object.defineProperty(workflow, 'canBeTerminated', {
		value: false
	});

	//

	// $pauseLiveUpdates = true;
	$workflowRun = { ...$workflowRun, workflow };
	$fullEventHistory = toEventHistory(eventHistory);
	$currentEventHistory = toEventHistory(eventHistory);

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
