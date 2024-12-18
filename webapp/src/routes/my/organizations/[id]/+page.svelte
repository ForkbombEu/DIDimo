<script lang="ts">
	import OrganizationLayout from '@/organizations/components/organizationLayout.svelte';

	const { data } = $props();
	const { organization } = $derived(data);
</script>

<OrganizationLayout org={organization}>
	<pre>{JSON.stringify(organization, null, 2)}</pre>
</OrganizationLayout>
