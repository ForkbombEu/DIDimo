<script lang="ts" context="module">
	export type SidebarLayoutMode = 'default' | 'drawer';
</script>

<script lang="ts">
	import { getUIShellContext } from './UiShell.svelte';
	import SidebarDefaultContainer from './SidebarDefaultContainer.svelte';
	import SidebarDrawerContainer from './SidebarDrawerContainer.svelte';
	import { page } from '$app/stores';
	import { Sidebar } from 'flowbite-svelte';

	const { sidebarLayoutMode } = getUIShellContext();

	export let width = 'w-60';

	$: comp = $sidebarLayoutMode == 'drawer' ? SidebarDrawerContainer : SidebarDefaultContainer;
	$: activeUrl = $page.url.pathname;
	export let activeClass =
		'flex items-center p-2 text-base font-normal text-gray-900 bg-gray-200 dark:bg-gray-700 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700';
</script>

<svelte:component this={comp} {width}>
	<Sidebar
		class="flex flex-col grow bg-white border-r"
		asideClass={width}
		{activeUrl}
		{activeClass}
	>
		<div class="flex flex-col overflow-hidden grow">
			{#if $sidebarLayoutMode == 'drawer'}
				<div class="shrink-0">
					<slot name="top" />
				</div>
			{/if}
			<div class="h-0 grow overflow-y-auto">
				<slot />
			</div>
			<div class="shrink-0">
				<slot name="bottom" />
			</div>
		</div>
	</Sidebar>
</svelte:component>
