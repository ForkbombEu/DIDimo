<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import type { PBResponse } from '$lib/utils/types';
	type RecordGeneric = $$Generic<PBResponse>;

	export let value: unknown;
</script>

<div class="bg-gray-200 w-fit py-1 px-2 rounded-full text-xs">
	{value}
</div>
