<script lang="ts">
	import OrganizationBreadcrumbs from './organizationBreadcrumbs.svelte';
	import { PageTop, PageContent } from '@/components/layout';
	import OrganizationTabs from './organizationTabs.svelte';
	import type { OrganizationsResponse } from '@/pocketbase/types';
	import T from '@/components/ui-custom/t.svelte';

	interface Props {
		org: OrganizationsResponse;
		children?: import('svelte').Snippet;
	}

	let { org, children }: Props = $props();
</script>

<PageTop>
	<OrganizationBreadcrumbs />
	<T tag="h4">{org.name}</T>
	{#snippet bottom()}
		<OrganizationTabs organizationId={org.id} />
	{/snippet}
</PageTop>

<PageContent>
	{@render children?.()}
</PageContent>
