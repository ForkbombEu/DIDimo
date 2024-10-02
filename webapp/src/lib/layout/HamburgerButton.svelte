<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script>
	import { Button } from 'flowbite-svelte';
	import { getUIShellContext } from './UiShell.svelte';
	import { Bars3 } from 'svelte-heros-v2';

	const { toggleSidebar } = getUIShellContext();
</script>

<Button color="alternative" class="!p-1.5" on:click={toggleSidebar}>
	<Bars3 size="24" />
</Button>
