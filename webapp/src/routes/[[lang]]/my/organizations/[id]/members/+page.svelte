<script lang="ts">
	import { CollectionManager } from '$lib/collectionManager';
	import CollectionTable from '$lib/collectionManager/ui/collectionTable.svelte';
	import CreateRecord from '$lib/collectionManager/ui/recordActions/createRecord.svelte';
	import {
		Collections,
		type OrgAuthorizationsResponse,
		type OrgRolesResponse,
		type UsersResponse
	} from '$lib/pocketbase/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { Heading, P } from 'flowbite-svelte';

	export let data;
	$: organization = data.organization;

	const recordType =
		createTypeProp<OrgAuthorizationsResponse<{ user: UsersResponse; role: OrgRolesResponse }>>();
</script>

<CollectionManager
	{recordType}
	collection={Collections.OrgAuthorizations}
	formSettings={{
		hide: { organization: organization.id },
		relations: {
			role: { inputMode: 'select', displayFields: ['name'] },
			user: { displayFields: ['email'] }
		}
	}}
	editFormSettings={{
		exclude: ['user']
	}}
	initialQueryParams={{ expand: 'user,role', filter: `organization.id="${organization.id}"` }}
	let:records
>
	<div class="flex justify-between mb-4 items-center">
		<Heading tag="h6">Manage members and roles</Heading>
		<div class="shrink-0">
			<CreateRecord />
		</div>
	</div>
	<CollectionTable {records} hideActions={['select', 'share']} let:record>
		<P class="space-x-4">
			<span>{record.expand?.user.email}</span>
			<span>{record.expand?.role.name}</span>
		</P>
	</CollectionTable>
</CollectionManager>
