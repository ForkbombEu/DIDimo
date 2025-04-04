<script>
	import BackButton from '$lib/layout/back-button.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';

	let { data } = $props();
</script>

<PageTop contentClass="!space-y-0">
	<BackButton href="/my/workflows">Back to workflows</BackButton>
	<T tag="h1">Workflow {data.workflow_id}</T>
</PageTop>

<pre>{JSON.stringify(data, null, 2)}</pre>
