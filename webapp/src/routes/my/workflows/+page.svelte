<script lang="ts">
	import BackButton from '$lib/layout/back-button.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import A from '@/components/ui-custom/a.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import * as Table from '@/components/ui/table/index.js';
	import { toWorkflowExecution } from '@forkbombeu/temporal-ui';

	let { data } = $props();
	const { executions } = $derived(data);
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
				<Table.Head class="text-right">Start</Table.Head>
				<Table.Head class="text-right">End</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each executions as workflow (workflow.execution.runId)}
				{@const path = `/my/workflows/${workflow.execution.workflowId}/${workflow.execution.run_id}`}
				<Table.Row>
					<Table.Cell>{workflow.status}</Table.Cell>
					<Table.Cell class="font-medium">
						<A href={path}>{workflow.execution.workflowId}</A>
					</Table.Cell>
					<Table.Cell>{workflow.type.name}</Table.Cell>
					<Table.Cell class="text-right">{workflow.startTime}</Table.Cell>
					<Table.Cell class="text-right">{workflow.endTime}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</PageContent>
