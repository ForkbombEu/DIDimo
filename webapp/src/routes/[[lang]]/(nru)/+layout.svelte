<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Register } from 'flowbite-svelte-blocks';
	import { appTitle } from '$lib/strings';
	import BaseLanguageSwitcher from '$lib/i18n/baseLanguageSwitcher.svelte';
</script>

<div class="grid h-screen place-items-center">
	<div class="w-2/5">
		<Register href="/">
			<svelte:fragment slot="top">
				<img class="mr-2 h-10" src="/logo.svg" alt={`${appTitle} Logo`} />
			</svelte:fragment>
			<div class="space-y-4 p-6 sm:p-8 md:space-y-6">
				<slot />
			</div>
		</Register>
	</div>
</div>

<div class="fixed right-4 top-4 !m-0">
	<BaseLanguageSwitcher />
</div>
