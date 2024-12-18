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
