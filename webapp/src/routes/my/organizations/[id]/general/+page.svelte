<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { Collections, type OrganizationsResponse } from '$lib/pocketbase/types';
	import { RecordForm } from '$lib/recordForm';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { Heading } from 'flowbite-svelte';

	//

	export let data;
	$: organization = data.organization;

	const recordType = createTypeProp<OrganizationsResponse>();
</script>

<Heading tag="h6" class="mb-6">Manage your organization public info</Heading>
<RecordForm
	{recordType}
	collection={Collections.Organizations}
	recordId={organization.id}
	initialData={organization}
	submitButtonText="Save changes"
	on:success={invalidateAll}
/>
