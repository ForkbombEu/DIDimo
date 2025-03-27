<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables } from './logic';

	//

	let { data } = $props();

	let d = $state<Record<string, unknown>>();
</script>

<!--  -->

{#if !d}
	<SelectTestForm
		standards={data.standardsAndTestSuites}
		onSelectTests={(standardId, tests) => {
			getVariables(standardId, tests).then((res) => {
				d = res;
			});
		}}
	/>
{:else}
	<pre>{JSON.stringify(d, null, 2)}</pre>
{/if}
