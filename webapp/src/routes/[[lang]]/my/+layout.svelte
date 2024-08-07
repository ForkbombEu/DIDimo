<script lang="ts">
	import { currentUser } from '$lib/pocketbase';
	import { goto } from '$lib/i18n';

	import {
		UIShell,
		Sidebar,
		Topbar,
		HamburgerButton,
		AvatarMenu,
		Logo,
		SidebarLinks,
		MainContent,
		SidebarCloseButton
	} from '$lib/layout';
	import { DropdownDivider, DropdownHeader, DropdownItem } from 'flowbite-svelte';
	import DIDButton from '$lib/components/DIDButton.svelte';
	import { Fire } from 'svelte-heros-v2';

	import { links } from './sidebarLinks';

	let sidebarLayoutBreakpoint = 1024;
</script>

<UIShell {sidebarLayoutBreakpoint}>
	<Topbar slot="top" let:sidebarLayoutMode>
		<svelte:fragment slot="left">
			<div class="flex space-x-2">
				<HamburgerButton />
				{#if sidebarLayoutMode == 'default'}
					<Logo />
				{/if}
			</div>
		</svelte:fragment>
		<svelte:fragment slot="center">
			<div class="flex items-center">
				<div>
					<span class="whitespace-nowrap">
						Hello, <span class="text-primary-600 font-semibold">{$currentUser?.email}</span>
					</span>
				</div>
				<div class="shrink-0">
					<DIDButton />
				</div>
			</div>
		</svelte:fragment>
		<svelte:fragment slot="right">
			<AvatarMenu>
				<DropdownHeader>
					<span class="block truncate text-xs font-medium text-gray-500">
						{$currentUser?.email}
					</span>
				</DropdownHeader>
				<DropdownItem href="/my/profile">My profile</DropdownItem>
				<DropdownDivider />
				<DropdownItem href="/pro" class="flex items-center">
					<Fire class="mr-2 w-5 text-red-500" /> Go Pro</DropdownItem
				>
				<DropdownDivider />
				<DropdownItem on:click={() => goto('/logout')} class="text-primary-600">
					Sign out
				</DropdownItem>
			</AvatarMenu>
		</svelte:fragment>
	</Topbar>

	<Sidebar>
		<svelte:fragment slot="top">
			<div class="flex items-center p-3">
				<Logo />
				<SidebarCloseButton />
			</div>
		</svelte:fragment>
		<div class="p-3">
			<SidebarLinks {links} />
		</div>
		<svelte:fragment slot="bottom">bottom</svelte:fragment>
	</Sidebar>

	<MainContent>
		<slot />
	</MainContent>
</UIShell>
