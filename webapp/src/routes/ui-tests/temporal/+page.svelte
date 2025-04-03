<script lang="ts">
	import { workflowResponse, eventHistory } from './data';

	import TemporalI18nProvider from './temporal-i18n-provider.svelte';
	import TemporalWorkflow from './temporal-workflow.svelte';
</script>

<TemporalI18nProvider>
	<TemporalWorkflow {workflowResponse} {eventHistory} />
</TemporalI18nProvider>
