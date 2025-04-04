<script lang="ts">
	import BackButton from '$lib/layout/back-button.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import A from '@/components/ui-custom/a.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import * as Table from '@/components/ui/table/index.js';
	import { toWorkflowExecution } from '@forkbombeu/temporal-ui';

	let { data } = $props();
	const executions = $derived(data.executions);

	/**
	 * TODO - This has to work
	 *
	 * The error is:
	 * Cannot read properties of undefined (reading 'searchAttributes')
	 *
	 * - `status` must be converted from code to string
	 * - `execution_time` should be converted
	 */
	const parsedExecutions = $derived(
		executions.map((e) => toWorkflowExecution({ workflowExecutionInfo: e }))
	);
	$inspect(parsedExecutions);
	/**
	 * `e` should have shape:
	 *
	 * export type WorkflowExecutionInfo = Replace<WorkflowExeuctionWithAssignedBuildId, {
	 * 	status: WorkflowExecutionStatus | WorkflowStatus;
	 * 	stateTransitionCount: string;
	 * 	startTime: string;
	 * 	closeTime: string;
	 * 	executionTime: string;
	 * 	historySizeBytes: string;
	 * 	historyLength: string;
	 * 	searchAttributes?: WorkflowSearchAttributes;
	 * 	memo?: Memo;
	 * }>;
	 */
</script>

<PageTop contentClass="!space-y-0">
	<BackButton href="/my">Back to dashboard</BackButton>
	<T tag="h1">Workflows</T>
</PageTop>

<PageContent class="bg-secondary grow">
	<Table.Root class="rounded-lg bg-white">
		<Table.Header>
			<Table.Row>
				<Table.Head>Status</Table.Head>
				<Table.Head>Workflow ID</Table.Head>
				<Table.Head>Type</Table.Head>
				<Table.Head class="text-right">Execution time</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each executions as workflow (workflow.execution.workflow_id)}
				{@const path = `/my/workflows/${workflow.execution.workflow_id}/${workflow.execution.run_id}`}
				<Table.Row>
					<Table.Cell>{workflow.status}</Table.Cell>
					<Table.Cell class="font-medium">
						<A href={path}>{workflow.execution.workflow_id}</A>
					</Table.Cell>
					<Table.Cell>{workflow.type.name}</Table.Cell>
					<Table.Cell class="text-right">{workflow.execution_time.seconds}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</PageContent>
