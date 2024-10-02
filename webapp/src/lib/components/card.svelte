<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	let className = '';
	export { className as class };
</script>

<div class={`bg-white rounded-md shadow-md dark:bg-gray-800 dark:text-white ${className}`}>
	<slot />
</div>
