<script lang="ts">
	import { WorkflowsWithNewSearch, workflows as workflowsStore } from '@forkbombeu/temporal-ui';
	import { page } from '$app/stores';
	import { getAllContexts, getContext } from 'svelte';

	// $page.params = { ...$page.params, namespace: '' };
	// const { workflows }: Props = $props();
	console.log($page.params.namespace);
</script>

{#if $page.params.namespace}
	<WorkflowsWithNewSearch></WorkflowsWithNewSearch>
{/if}
