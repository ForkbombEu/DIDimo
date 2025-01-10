<script lang="ts">
	import BaseLayout from '$lib/layout/baseLayout.svelte';
	// import OrganizationSwitcher from '$lib/layout/organizationSwitcher.svelte';
	// import { Button } from '@/components/ui/button';
	// import { m } from '@/i18n';
	// import { currentUser } from '@/pocketbase';
	// import { PocketbaseQuery } from '@/pocketbase/query';
	import type { Snippet } from 'svelte';

	//

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();

	//

	// const authorizationsQuery = new PocketbaseQuery('orgAuthorizations', {
	// 	filter: `user = "${$currentUser!.id}"`,
	// 	expand: ['organization', 'role']
	// });

	// const organizationsPromise = authorizationsQuery.getFullList().then((authorizations) => {
	// 	return authorizations.map((authorization) => authorization.expand?.organization!);
	// });
</script>

<BaseLayout>
	<!-- <div class="border-y">
		<div class="mx-auto flex max-w-screen-xl items-center px-2 py-2">
			{#await organizationsPromise then organizations}
				<OrganizationSwitcher {organizations} />
			{/await}
			<Button variant="link" href="/my/profile">{m.My_profile()}</Button>
			<Button variant="link" href="/my/organizations">{m.organizations()}</Button>
		</div>
	</div> -->
	{@render children?.()}
</BaseLayout>
