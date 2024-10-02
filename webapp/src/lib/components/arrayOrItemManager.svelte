<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import IconButton from './iconButton.svelte';
	import { Listgroup, ListgroupItem } from 'flowbite-svelte';
	import ListgroupItemSlotted from './listgroupItemSlotted.svelte';

	//

	type T = $$Generic;
	export let value: T | T[] | undefined;

	//

	let tempItems: T[] = [];
	$: createTempItems(value);

	function createTempItems(v: typeof value) {
		if (Array.isArray(v)) tempItems = v;
		else if (v) tempItems = [v];
		else tempItems = [];
	}

	function removeItem(item: T) {
		if (Array.isArray(value)) value = value.filter((i) => i !== item);
		else value = undefined;
	}
</script>

{#if tempItems.length > 0}
	<Listgroup>
		{#each tempItems as item (item)}
			<ListgroupItemSlotted>
				<svelte:fragment slot="left">
					<slot {item} />
				</svelte:fragment>
				<svelte:fragment slot="right">
					<slot name="right" {item} />
					<IconButton
						on:click={() => {
							removeItem(item);
						}}
					/>
				</svelte:fragment>
			</ListgroupItemSlotted>
		{/each}
	</Listgroup>
{/if}
