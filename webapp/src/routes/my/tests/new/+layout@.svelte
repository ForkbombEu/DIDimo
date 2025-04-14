<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script>
	import BaseLayout from '$lib/layout/baseLayout.svelte';

	let { children } = $props();
</script>

<BaseLayout>
	{@render children()}
</BaseLayout>
