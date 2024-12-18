<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';

	let { data } = $props();
	const did = JSON.stringify(data.did, null, 2);
</script>

<div class="space-y-4 p-4">
	<T tag="h3">Your DID</T>
	<pre class="rounded-lg border border-slate-200 bg-white p-4">{did}</pre>
</div>
