<script lang="ts">
	import { createRecordLabel } from '$lib/components/records/utils';

	import type { z } from 'zod';
	import type { FormPath, FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	import RecordsManager, {
		type RecordsManagerOptions
	} from '$lib/components/records/recordsManager.svelte';
	import type { Collections } from '$lib/pocketbase/types';
	import FieldWrapper from './fieldParts/fieldWrapper.svelte';
	import { createTypeProp } from '$lib/utils/typeProp';
	import type { HelpTextOption, LabelOption } from './types';
	import type { PBResponse } from '$lib/utils/types';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	type T = $$Generic<AnyZodObject>;
	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>> | FormPath<z.infer<T>>;
	export let options: Partial<RecordsManagerOptions<RecordGeneric>> & LabelOption & HelpTextOption =
		{};
	export let collection: string | Collections;

	//

	const { validate } = superform;
	const { value } = formFieldProxy(superform, field as string);
	if (!options.name) options.name = field;

	$: if (Array.isArray($value) && $value.length > 0) validate(field, $value as any);
</script>

<FieldWrapper {field} label={options.label}>
	<slot slot="labelRight" name="labelRight" />
	<RecordsManager {recordType} {collection} bind:value={$value} options={{ ...options }} let:record>
		<slot {record}>
			{createRecordLabel(record, options.displayFields)}
		</slot>
	</RecordsManager>
</FieldWrapper>
