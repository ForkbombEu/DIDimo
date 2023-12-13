<script lang="ts">
	import { CloseButton } from 'flowbite-svelte';
	import { getUIShellContext } from './UiShell.svelte';

	const { isSidebarHidden } = getUIShellContext();
</script>

<CloseButton on:click={() => ($isSidebarHidden = true)} />
