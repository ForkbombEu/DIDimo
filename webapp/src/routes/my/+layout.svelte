<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import BaseLayout from '$lib/layout/baseLayout.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import NavigationTabs from '@/components/ui-custom/navigationTabs.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import { currentUser } from '@/pocketbase';
	// import OrganizationSwitcher from '$lib/layout/organizationSwitcher.svelte';
	// import Button from '@/components/ui-custom/button.svelte'
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
	<PageTop>
		<T tag="h1">{m.hello_user({ username: $currentUser?.name! })}</T>

		<NavigationTabs
			tabs={[
				{ title: 'Services', href: '/my/services' },
				{ title: 'Test runs', href: '/my/tests/runs' },
				{ title: 'Profile', href: '/my/profile' }
			]}
		/>
	</PageTop>

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
