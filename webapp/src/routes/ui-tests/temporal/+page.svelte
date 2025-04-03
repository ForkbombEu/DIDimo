<script lang="ts">
	import { workflowResponse, eventHistory } from './data';

	import TemporalContextProvider from './temporal-context-provider.svelte';
	import TemporalWorkflow from './temporal-workflow.svelte';
	import TemporalWorkflowList from './temporal-workflow-list.svelte';
</script>

<TemporalContextProvider>
	<TemporalWorkflowList />
	<hr />
	<TemporalWorkflow {workflowResponse} {eventHistory} />
</TemporalContextProvider>
