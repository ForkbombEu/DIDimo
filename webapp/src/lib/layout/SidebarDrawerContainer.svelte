<script lang="ts">
	import Drawer from '$lib/components/drawer.svelte';
	import { getUIShellContext } from './UiShell.svelte';

	const { isSidebarHidden } = getUIShellContext();

	export let placement: 'left' | 'right' = 'left';
	export let width = 'w-40';
	export let backdrop = true;
	export let closeOnClickOutside = true;
</script>

<Drawer
	bind:hidden={$isSidebarHidden}
	{placement}
	{width}
	{backdrop}
	{closeOnClickOutside}
	hideTopbar
>
	<slot />
</Drawer>
