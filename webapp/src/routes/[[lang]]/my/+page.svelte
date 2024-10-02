<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script>
	import { A } from 'flowbite-svelte';
	import { featureFlags } from '$lib/features';
</script>

<div class="flex flex-col p-4">
	<A href="/my/profile">My profile</A>
	{#if $featureFlags.ORGANIZATIONS}
		<A href="/my/organizations">My organizations</A>
	{/if}
</div>
