<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	export let text: string | undefined = undefined;
</script>

{#if text}
	<p class="text-xs text-gray-400">{text}</p>
{/if}
