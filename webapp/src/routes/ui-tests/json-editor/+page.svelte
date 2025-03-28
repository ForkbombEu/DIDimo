<script lang="ts">
	import SelectTestForm from './select-test-form.svelte';
	import { getVariables, type FieldsResponse } from './logic';
	import TestsDataForm from './tests-data-form.svelte';

	//

	let { data } = $props();

	let d = $state<FieldsResponse>();

	const x = [
		'iso_mdl:pre_registered:request_uri_signed:direct_post.jwt.json',
		'iso_mdl:pre_registered:request_uri_signed:w3c_dc_api.json',
		'iso_mdl:pre_registered:request_uri_signed:w3c_dc_api.jwt.json',
		'iso_mdl:pre_registered:request_uri_unsigned:direct_post.json',
		'iso_mdl:pre_registered:request_uri_unsigned:direct_post.jwt.json',
		'iso_mdl:pre_registered:request_uri_unsigned:w3c_dc_api.json',
		'iso_mdl:pre_registered:request_uri_unsigned:w3c_dc_api.jwt.json',
		'iso_mdl:redirect_uri:request_uri_signed:direct_post.json',
		'iso_mdl:redirect_uri:request_uri_signed:direct_post.jwt.json',
		'iso_mdl:redirect_uri:request_uri_signed:w3c_dc_api.json',
		'sd_jwt_vc:x509_san_dns:request_uri_signed:direct_post.jwt.json',
		'sd_jwt_vc:x509_san_dns:request_uri_signed:w3c_dc_api.json',
		'sd_jwt_vc:x509_san_dns:request_uri_signed:w3c_dc_api.jwt.json',
		'sd_jwt_vc:x509_san_dns:request_uri_unsigned:direct_post.json'
	];

	getVariables('', x).then((res) => {
		d = res;
		// console.log(d);
	});
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
	<TestsDataForm data={d} />
{/if}
