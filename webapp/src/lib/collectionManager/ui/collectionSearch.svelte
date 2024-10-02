<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { createTypeProp } from '$lib/utils/typeProp';
	import type { PBResponse, StringKeys } from '$lib/utils/types';

	import { Button, ButtonGroup, Dropdown, Search, Checkbox } from 'flowbite-svelte';
	import { getRecordsManagerContext } from '../collectionManager.svelte';
	import { getCollectionSchema } from '$lib/pocketbase/schema';
	import { ChevronDown } from 'svelte-heros-v2';

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	export let searchableFields: StringKeys<RecordGeneric>[] = [];
	export let placeholder: string = 'Search...';
	export let allFieldsCaption: string = 'All fields';
	export let filtersActiveCaption: string = 'Filters active';

	const { dataManager, collection } = getRecordsManagerContext();
	const { queryParams } = dataManager;
	const allFields = getCollectionSchema(collection)?.schema?.map((field) => field.name) ?? [];
	let selected: string[] = [allFieldsCaption];
	let fields: { value: string; name: string }[] = [allFieldsCaption]
		.concat(searchableFields || allFields)
		.map((field) => ({
			value: field,
			name: field
		}));
	let queryString = '';
	const handleSearch = () => {
		if (selected.includes(allFieldsCaption)) {
			$queryParams.filter = (searchableFields || allFields)
				.reduce((acc: string[], field) => {
					acc.push(`${field} ~ '${queryString}'`);
					return acc;
				}, [])
				.join(' || ');
		} else {
			$queryParams.filter = selected
				.reduce((acc: string[], field) => {
					acc.push(`${field} ~ '${queryString}'`);
					return acc;
				}, [])
				.join(' || ');
		}
	};
	$: {
		if (selected.slice(-1)[0] === allFieldsCaption) {
			selected = [allFieldsCaption];
		} else if (selected.length > 1 && selected.includes(allFieldsCaption)) {
			selected = selected.filter((item) => item !== allFieldsCaption);
		}
	}
</script>

<form class="flex gap-2">
	<ButtonGroup class="w-full z-20">
		<Button
			color="none"
			class="flex-shrink-0 text-gray-900 bg-gray-100 border border-gray-300 dark:border-gray-700 dark:text-white hover:bg-gray-200 focus:ring-gray-300 dark:bg-gray-600 dark:hover:bg-gray-700 dark:focus:ring-gray-800"
		>
			<ChevronDown class="w-4 h-4 mr-2" />

			{selected.includes(allFieldsCaption) ? allFieldsCaption : filtersActiveCaption}
		</Button>
		<Dropdown class="p-3 space-y-1">
			{#each fields as field}
				<li class="!w-fit font-semibold px-2">
					<Checkbox
						checked={selected.includes(field.value)}
						on:change={() => {
							if (selected.includes(field.value)) {
								console.log('deselected', selected, field);
								selected = selected.filter((item) => item !== field.value);
							} else {
								console.log('selected', selected, field);
								selected = selected.concat(field.value);
							}
						}}
					>
						{field.name}
					</Checkbox>
				</li>
			{/each}
		</Dropdown>
		<Search
			size="lg"
			bind:value={queryString}
			on:input={() => {
				handleSearch();
			}}
			{placeholder}
		/>
	</ButtonGroup>
</form>
