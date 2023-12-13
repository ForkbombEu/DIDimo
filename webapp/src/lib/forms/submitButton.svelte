<script lang="ts">
	import { Button } from 'flowbite-svelte';
	import { formHasErrors, getFormContext } from './form.svelte';
	import type { ComponentProps } from 'svelte';

	export let id = 'submit';
	export let color: ComponentProps<Button>['color'] = 'primary';

	const { superform } = getFormContext();
	const { allErrors } = superform;

	$: hasErrors = formHasErrors($allErrors);
</script>

<Button {id} {color} type="submit" disabled={hasErrors}>
	<slot />
</Button>
