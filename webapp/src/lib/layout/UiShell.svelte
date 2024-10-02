<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script context="module" lang="ts">
	import type { Writable } from 'svelte/store';
	import { getContext } from 'svelte';
	import type { SidebarLayoutMode } from './Sidebar.svelte';

	export const UI_SHELL_KEY = Symbol('usk');

	export type UIShellContext = {
		isSidebarHidden: Writable<boolean>;
		toggleSidebar: () => void;
		sidebarLayoutMode: Writable<SidebarLayoutMode>;
	};

	export function getUIShellContext(): UIShellContext {
		return getContext(UI_SHELL_KEY);
	}
</script>

<script lang="ts">
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';

	export let sidebarLayoutBreakpoint: number | undefined = undefined;

	let windowWidth: number;

	const isSidebarHidden = writable(false);
	const sidebarLayoutMode = writable<SidebarLayoutMode>('default');

	$: if (sidebarLayoutBreakpoint && windowWidth && windowWidth >= sidebarLayoutBreakpoint) {
		$isSidebarHidden = false;
		$sidebarLayoutMode = 'default';
	} else {
		$isSidebarHidden = true;
		$sidebarLayoutMode = 'drawer';
	}

	const toggleSidebar = () => {
		$isSidebarHidden = !$isSidebarHidden;
	};

	setContext<UIShellContext>(UI_SHELL_KEY, {
		isSidebarHidden,
		toggleSidebar,
		sidebarLayoutMode
	});
</script>

<svelte:window bind:innerWidth={windowWidth} />

<div class="w-screen h-screen overflow-hidden flex flex-col">
	<div class="shrink-0">
		<slot name="top" sidebarLayoutMode={$sidebarLayoutMode} />
	</div>
	<div class="h-0 grow flex items-stretch">
		<slot sidebarLayoutMode={$sidebarLayoutMode} />
	</div>
</div>
