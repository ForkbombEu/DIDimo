<script lang="ts">
	import BackButton from '$lib/layout/back-button.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import A from '@/components/ui-custom/a.svelte';
	import T from '@/components/ui-custom/t.svelte';
	// import { toWorkflowExecution } from '@forkbombeu/temporal-ui';

	let { data } = $props();
	const executions = $derived(data.executions);

	// const executions = $derived(toWorkflowExecutions(data.executions));
	// const executions = $derived(data.executions.map((e) => toWorkflowExecution(e)));
</script>

<PageTop contentClass="!space-y-0">
	<BackButton href="/my">Back to dashboard</BackButton>
	<T tag="h1">Workflows</T>
</PageTop>

<!-- <pre>{JSON.stringify(executions, null, 2)}</pre> -->

{#each executions as workflow}
	{@const path = `/my/workflows/${workflow.execution.workflow_id}/${workflow.execution.run_id}`}
	<div class="flex items-center gap-2">
		<T>{workflow.execution.run_id}</T>
		<T><A href={path}>{workflow.execution.workflow_id}</A></T>
		<T>{workflow.type.name}</T>
		<T>{workflow.status}</T>
		<T>{workflow.execution_time.seconds}</T>
		<T>{workflow.execution_time.nanos}</T>
	</div>
{/each}
