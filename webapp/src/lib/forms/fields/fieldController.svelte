<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;

	const { value, errors, constraints } = formFieldProxy(superform, field);

	function updateValue(newValue: typeof $value) {
		value.set(newValue);
	}
</script>

<slot value={$value} errors={$errors} constraints={$constraints} {updateValue} {field} />
