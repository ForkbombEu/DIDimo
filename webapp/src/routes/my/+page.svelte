<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import PageTop from '$lib/layout/pageTop.svelte';
	import { WelcomeSession } from '@/auth/welcome';
	import { CollectionManager } from '@/collections-components/index.js';
	import { RecordCard } from '@/collections-components/manager';
	import Button from '@/components/ui-custom/button.svelte';
	import Icon from '@/components/ui-custom/icon.svelte';
	import NavigationTabs from '@/components/ui-custom/navigationTabs.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Separator } from '@/components/ui/separator';
	import { m } from '@/i18n/index.js';
	import { currentUser } from '@/pocketbase';
	import { Sparkle, Workflow } from 'lucide-svelte';

	if (WelcomeSession.isActive()) WelcomeSession.end();

	let { data } = $props();
</script>

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

<div class="mx-auto flex w-full max-w-screen-xl flex-col space-y-8 p-4">
	<Separator />

	<div class="flex items-center justify-between">
		<div class="space-y-2"></div>

		<div class="flex items-center gap-2">
			<Button href="/my/tests/runs" variant="outline">
				<Icon src={Workflow} />
				{m.View_test_runs()}
			</Button>
		</div>
	</div>

	<Separator />
</div>
