<script lang="ts">
	import { Button, Heading } from 'flowbite-svelte';

	export let data;
	const { organization } = data;
</script>

<pre>{JSON.stringify(organization, null, 2)}</pre>
