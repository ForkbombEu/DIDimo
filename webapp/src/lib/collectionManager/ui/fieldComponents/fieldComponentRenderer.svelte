<script lang="ts" context="module">
	import type { PBResponse } from '$lib/utils/types';
	import type { SvelteComponent } from 'svelte';

	export type FieldComponent<RecordGeneric extends PBResponse = PBResponse> =
		typeof SvelteComponent<{
			value?: unknown;
			record?: RecordGeneric;
		}>;
</script>

<script lang="ts">
	type RecordGeneric = $$Generic<PBResponse>;

	export let record: RecordGeneric;
	export let field: string;
	export let component: FieldComponent | undefined;
	export let showLabel: boolean = false;

	$: value = record[field];
	$: text = showLabel ? `${field}: ${value}` : value;
</script>

{#if component}
	<svelte:component this={component} {record} {value} />
{:else}
	<div>{text}</div>
{/if}
