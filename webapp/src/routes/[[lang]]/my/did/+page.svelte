<script lang="ts">
	import { Heading } from 'flowbite-svelte';

	export let data;
	const did = JSON.stringify(data.did, null, 2);
</script>

<div class="p-4 space-y-4">
	<Heading tag="h3">Your DID</Heading>
	<pre class="bg-white p-4 border border-slate-200 rounded-lg">{did}</pre>
</div>
