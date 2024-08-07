<script lang="ts">
	import { getFormContext } from '../../form.svelte';
	import { formFieldProxy } from 'sveltekit-superforms/client';

	export let field: string;

	const { superform, showRequiredIndicator } = getFormContext();
	const { constraints } = formFieldProxy(superform, field);
</script>

{#if $constraints?.required && showRequiredIndicator}
	<span class="font-bold text-red-400">*</span>
{/if}
