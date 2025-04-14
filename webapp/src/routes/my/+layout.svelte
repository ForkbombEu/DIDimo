<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import BaseLayout from '$lib/layout/baseLayout.svelte';
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import NavigationTabs from '@/components/ui-custom/navigationTabs.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { m } from '@/i18n';
	import { currentUser } from '@/pocketbase';
	import { GlobeIcon, Home, Shapes, TestTubeDiagonalIcon, User } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	//

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<BaseLayout>
	<PageTop contentClass="!pb-0">
		<T tag="h1">{m.hello_user({ username: $currentUser?.name! })}</T>

		<NavigationTabs
			tabs={[
				{ title: 'Services and Products', href: '/my/services-and-products', icon: Shapes },
				{ title: 'Test runs', href: '/my/tests/runs', icon: TestTubeDiagonalIcon },
				{ title: 'Organization page', href: '/my/organization-page', icon: GlobeIcon },
				{ title: 'Profile', href: '/my/profile', icon: User }
			]}
		/>
	</PageTop>

	<PageContent class="grow bg-secondary">
		{@render children?.()}
	</PageContent>
</BaseLayout>
