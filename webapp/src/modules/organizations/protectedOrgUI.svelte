<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { verifyUserRole } from './verify-authorizations';
	import type { OrgRole } from '.';

	interface Props {
		orgId: string;
		roles: OrgRole[];
		children?: import('svelte').Snippet;
	}

	let { orgId, roles, children }: Props = $props();
</script>

{#await verifyUserRole(orgId, roles) then response}
	{#if response.hasRole}
		{@render children?.()}
	{/if}
{/await}
