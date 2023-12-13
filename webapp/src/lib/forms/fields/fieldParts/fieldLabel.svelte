<script lang="ts">
	import { getFormContext } from '../../form.svelte';
	import { formFieldProxy } from 'sveltekit-superforms/client';
	import { fieldHasErrors } from './fieldError.svelte';
	import { Label } from 'flowbite-svelte';
	import FieldRequiredIndicator from './fieldRequiredIndicator.svelte';
	import { capitalizeFirstLetter } from '$lib/utils/strings';

	export let field: string;
	export let text = '';
	if (!text) text = capitalizeFirstLetter(field);

	const { superform } = getFormContext();
	const { errors } = formFieldProxy(superform, field);
</script>

<Label for={field} class="label" color={fieldHasErrors($errors) ? 'red' : 'gray'}>
	<span>{text}</span>
	<FieldRequiredIndicator {field} />
</Label>
