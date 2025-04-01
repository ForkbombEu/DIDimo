<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables, type FieldsResponse } from './logic';
	import TestsDataForm from './tests-data-form.svelte';

	//

	let { data } = $props();

	let d = $state<FieldsResponse>();
	let compositeTestId = $state('');
</script>

<!--  -->

{#if !d}
	<SelectTestForm
		standards={data.standardsAndTestSuites}
		onSelectTests={(standardId, tests) => {
			compositeTestId = standardId;
			getVariables(standardId, tests).then((res) => {
				d = res;
			});
		}}
	/>
{:else}
	<TestsDataForm data={d} testId={compositeTestId} />
{/if}
