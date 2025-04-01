<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables, type FieldsResponse } from './logic';
	import TestsDataForm from './tests-data-form.svelte';
	import T from '@/components/ui-custom/t.svelte';

	//

	let { data } = $props();

	let d = $state<FieldsResponse>();
	let compositeTestId = $state('');
</script>

<!--  -->

<div class="p-8">
	<T tag="h1">Compliance tests</T>
</div>

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
