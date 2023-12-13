<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	import RecordsManager, {
		type RecordsManagerOptions
	} from '$lib/components/records/recordsManager.svelte';
	import type { Collections } from '$lib/pocketbase/types';
	import FieldWrapper from './fieldParts/fieldWrapper.svelte';
	import { createTypeProp } from '$lib/utils/typeProp';
	import type { LabelOption } from './types';
	import type { PBResponse } from '$lib/utils/types';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	type T = $$Generic<AnyZodObject>;
	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: Partial<RecordsManagerOptions<RecordGeneric>> & LabelOption = {};
	export let collection: string | Collections;

	//

	const { validate } = superform;
	const { value } = formFieldProxy(superform, field as string);

	$: if (Array.isArray($value) && $value.length > 0) validate($value as any);
</script>

<FieldWrapper {field} label={options.label}>
	<RecordsManager {recordType} {collection} bind:value={$value} options={{ ...options }} />
</FieldWrapper>
