<script lang="ts">
	import BackButton from '$lib/layout/back-button.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';

	let { data } = $props();
	const { workflows } = $derived(data);
</script>

<PageTop contentClass="!space-y-2">
	<BackButton href="/my">Back to dashboard</BackButton>
	<T tag="h1">Workflows</T>
</PageTop>

<pre>{JSON.stringify(workflows, null, 2)}</pre>

<!-- {#each workflows as workflow}
	<div>
		<T>{workflow.type.name}</T>
		<T>{workflow.status}</T>
		<T>{workflow.execution_time.seconds}</T>
		<T>{workflow.execution_time.nanos}</T>
	</div>
{/each} -->
