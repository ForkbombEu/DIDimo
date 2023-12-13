<script lang="ts">
	import { verifyRole } from '.';
	import type { OrgRole } from './roles';

	export let orgId: string;
	export let roles: OrgRole[];
</script>

{#await verifyRole(orgId, roles) then _}
	<slot />
{/await}
