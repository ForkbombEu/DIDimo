<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { CloseButton } from 'flowbite-svelte';
	import { getUIShellContext } from './UiShell.svelte';

	const { isSidebarHidden } = getUIShellContext();
</script>

<CloseButton on:click={() => ($isSidebarHidden = true)} />
