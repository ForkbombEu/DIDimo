<script lang="ts">
	import { getFormContext } from '../../form.svelte';
	import { formFieldProxy } from 'sveltekit-superforms/client';

	export let field: string;

	const { superform, showRequiredIndicator } = getFormContext();
	const { constraints } = formFieldProxy(superform, field);
</script>

{#if $constraints?.required && showRequiredIndicator}
	<span class="text-red-400 font-bold">*</span>
{/if}
