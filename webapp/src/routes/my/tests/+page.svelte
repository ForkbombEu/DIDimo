<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables, type FieldsResponse } from './logic';
	import TestsDataForm from './tests-data-form.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { ArrowLeft } from 'lucide-svelte';
	import Button from '@/components/ui-custom/button.svelte';

	//

	let { data } = $props();

	let d = $state<FieldsResponse>();
	let compositeTestId = $state('');
</script>

<!--  -->

<div class="space-y-2 p-8">
	<Button href="/my" variant="link" class="p-0"><ArrowLeft size={16} /> Back to dashboard</Button>
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
