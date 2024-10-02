<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script>
	import { portal } from './utils/portal';

	export let target = 'body';
</script>

<div use:portal={target} class="fixed m-0 p-0 z-50">
	<slot />
</div>
