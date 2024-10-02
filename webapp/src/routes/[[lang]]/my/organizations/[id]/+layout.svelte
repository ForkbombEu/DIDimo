<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import { ProtectedOrgLayout, ProtectedOrgUI, OrgRoles } from '$lib/rbac';
	import { A, Button, ButtonGroup, Heading, Hr, P } from 'flowbite-svelte';

	export let data;
	$: organization = data.organization;
	$: avatarPath = pb.files.getUrl(organization, organization.avatar);

	const href = (path = '') => `/my/organizations/${organization.id}${path}`;
	const { OWNER, ADMIN } = OrgRoles;
</script>

<!--  -->

<A href="/my/organizations" class="mb-4">← My organizations</A>

<div class="flex justify-between items-center">
	<Heading tag="h3"><A href={href()}>{organization.name}</A></Heading>
	<ProtectedOrgUI orgId={organization.id} roles={[OWNER, ADMIN]}>
		<div class="flex justify-end items-center space-x-4">
			<P color="text-gray-600" size="sm">Settings</P>
			<ButtonGroup>
				<ProtectedOrgUI orgId={organization.id} roles={[OWNER]}>
					<Button href={href('/general')}>General</Button>
				</ProtectedOrgUI>
				<Button href={href('/members')}>Members</Button>
			</ButtonGroup>
		</div>
	</ProtectedOrgUI>
</div>

<Hr />

<ProtectedOrgLayout orgId={organization.id}>
	<slot />
</ProtectedOrgLayout>
