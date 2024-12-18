<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		class?: string;
		children?: Snippet;
	}

	let { class: className = '', children }: Props = $props();
</script>

<T class="text-card-foreground text-sm {className}">{@render children?.()}</T>
