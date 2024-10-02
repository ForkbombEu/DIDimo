<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { verifyRole } from '.';
	import type { OrgRole } from './roles';

	export let orgId: string;
	export let roles: OrgRole[];
</script>

{#await verifyRole(orgId, roles) then _}
	<slot />
{/await}
