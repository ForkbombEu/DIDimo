<script lang="ts">
	import FieldError from './fieldError.svelte';
	import FieldLabel from './fieldLabel.svelte';

	export let field: string;
	export let label: string | null | undefined = undefined;
</script>

<div class="space-y-2">
	{#if label}
		<FieldLabel {field} text={label} />
	{/if}
	<slot />
	<FieldError {field} />
</div>
