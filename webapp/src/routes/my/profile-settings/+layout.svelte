<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import NavSidebar from '@/components/ui-custom/navSidebar.svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
	};
	import type { LinkWithIcon } from '@/components/types';
	import { KeyIcon, LockIcon, UserIcon } from 'lucide-svelte';

	const { children }: Props = $props();

	const links: LinkWithIcon[] = [
		{
			href: '/my/profile-settings',
			title: 'Personal profile',
			icon: UserIcon
		},
		{
			href: '/my/profile-settings/public-keys',
			title: 'Public Keys',
			icon: KeyIcon
		},
		{
			href: '/my/profile-settings/webauthn',
			title: 'Webauthn',
			icon: LockIcon
		}
	];
</script>

<div class="flex w-full gap-8 p-8">
	<NavSidebar title="settings" {links}></NavSidebar>
	<div class="grow">
		{@render children()}
	</div>
</div>
