<script lang="ts">
	import type { RecordInputOptions } from './types';

	import { pb } from '$lib/pocketbase';
	import type { PBResponse } from '$lib/utils/types';
	import type { Collections } from '$lib/pocketbase/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { createRecordLabel, excludeStringArray, getCollectionFieldNames } from './utils';
	// @ts-ignore
	import Svelecte from 'svelecte';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	//

	export let collection: string | Collections;
	export let recordId: string | undefined = undefined;
	export let options: Partial<RecordInputOptions<RecordGeneric>> = {};

	let { displayFields = [], name = undefined, required = false, placeholder = undefined } = options;
	$: exclude = options.excludeIds ?? [];
	$: disabled = options.disabled ?? false;

	//

	const valueField = 'value';
	const labelField = 'label';

	type Option = {
		[valueField]: string;
		[labelField]: string;
	};

	function handleChange(e: CustomEvent<null | Option>) {
		recordId = e.detail?.[valueField];
	}

	//

	async function fetchOptions(text: string | undefined): Promise<Option[]> {
		if (!text) return [];

		const records = await pb.collection(collection).getFullList<RecordGeneric>({
			requestKey: null,
			filter: filterString(text)
		});

		return records.map((r) => {
			return {
				[valueField]: r.id,
				label: createRecordLabel(r, displayFields)
			};
		});
	}

	function filterString(text: string) {
		let baseString = filterSearchString(text);
		if (exclude.length > 0) {
			const excludeString = excludeStringArray('id', exclude);
			baseString = `${baseString} && ${excludeString}`;
		}
		return baseString;
	}

	function filterSearchString(text: string) {
		const searchString = getCollectionFieldNames(collection)
			.map((f) => `${f} ~ "${text}"`)
			.join(' || ');
		return `(${searchString})`;
	}

	//

	let record: RecordGeneric | undefined = undefined;
	$: loadRecord(recordId);

	let actualPlaceholder: string | undefined = placeholder;
	$: setPlaceholder(record);

	async function loadRecord(value: string | undefined) {
		if (value) record = await pb.collection(collection).getOne(value);
		else record = undefined;
	}

	function setPlaceholder(maybeRecord: typeof record) {
		if (maybeRecord) actualPlaceholder = createRecordLabel(maybeRecord, displayFields);
		else actualPlaceholder = placeholder;
	}
</script>

<Svelecte
	{name}
	{valueField}
	{labelField}
	placeholder={actualPlaceholder}
	fetch={fetchOptions}
	value={recordId}
	on:change={handleChange}
	{disabled}
	{required}
/>
