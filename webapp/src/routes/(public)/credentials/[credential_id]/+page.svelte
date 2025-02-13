<script>
	let { data } = $props();
	const { credential } = $derived(data);
</script>

<pre>{JSON.stringify(credential, null, 2)}</pre>
