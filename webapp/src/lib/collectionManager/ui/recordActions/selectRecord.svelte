<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { getRecordsManagerContext } from '../../collectionManager.svelte';
	import type { PBResponse } from '$lib/utils/types';
	import { Checkbox } from 'flowbite-svelte';

	type RecordGeneric = $$Generic<PBResponse>;
	export let record: RecordGeneric;

	let { selectionManager } = getRecordsManagerContext();
	let { selectedRecords } = selectionManager;
</script>

<Checkbox bind:group={$selectedRecords} value={record.id} name="select" />
