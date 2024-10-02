<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { pb } from '$lib/pocketbase';
	import { Button, Tooltip } from 'flowbite-svelte';
	import type { PBResponse } from '$lib/utils/types';
	import { ArrowDownTray } from 'svelte-heros-v2';

	type RecordGeneric = $$Generic<PBResponse>;
	export let record: RecordGeneric;
	export let value: string;

	let url = '';
	if (record) url = pb.files.getUrl(record, value);
</script>

{#if url}
	<Button href={url} target="_blank" class="!p-2" color="alternative">
		<ArrowDownTray size="20" />
	</Button>
	<Tooltip>Download: {value}</Tooltip>
{/if}
