<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables } from './logic';

	//

	let { data } = $props();
</script>

<!--  -->

<SelectTestForm
	standards={data.standardsAndTestSuites}
	onSelectTests={(standardId, tests) => {
		getVariables(standardId, tests).then((res) => {
			console.log(res);
		});
	}}
/>
