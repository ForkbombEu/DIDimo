<script lang="ts">
	import { OrgRoles } from '$lib/rbac';
	import { Heading, Button, A, P, Badge } from 'flowbite-svelte';
	import { ProtectedOrgUI } from '$lib/rbac';
	import { Plus } from 'svelte-heros-v2';
	import { c } from '$lib/utils/strings.js';

	export let data;
	let { authorizations } = data;

	const { ADMIN, OWNER } = OrgRoles;
</script>

<div class="flex justify-between items-center mb-6">
	<Heading tag="h5">Your organizations</Heading>
	<Button size="sm" color="alternative" class="!px-4 shrink-0" href="/my/organizations/create">
		<Plus size="20" />
		<span class="ml-1"> Create a new organization </span>
	</Button>
</div>

<div class="border rounded-lg divide-y">
	{#if authorizations.length == 0}
		<div class="p-4 text-gray-600">No organizations found. Create one!</div>
	{:else}
		{#each authorizations as a}
			{@const org = a.expand.organization}
			{@const role = a.expand.role}
			<div class="px-4 py-3 flex justify-between items-center">
				<div class="flex items-center space-x-4">
					<P href={`/my/organizations/${org.id}`}>{org.name}</P>
					{#if role.name == ADMIN || role.name == OWNER}
						<Badge large color="dark">{c(role.name)}</Badge>
					{/if}
				</div>
				<ProtectedOrgUI orgId={org.id} roles={[ADMIN, OWNER]}>
					<Button
						data-testid={`${org.name} link`}
						size="sm"
						color="alternative"
						href={`/my/organizations/${org.id}`}
					>
						Settings
					</Button>
				</ProtectedOrgUI>
			</div>
		{/each}
	{/if}
</div>
