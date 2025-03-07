<script lang="ts">
	import type { Snippet } from 'svelte';
	import T from '@/components/ui-custom/t.svelte';

	type Props = {
		label?: string;
		children: Snippet;
	};

	let { label, children }: Props = $props();
</script>

<div class="space-y-1">
	{#if label}
		<T class="text-sm">{label}:</T>
	{/if}

	<div class="w-fit rounded-sm border border-slate-400 bg-white px-2 py-1">
		{@render children()}
	</div>
</div>
