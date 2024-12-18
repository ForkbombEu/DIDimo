<script lang="ts">
	import type { Snippet } from 'svelte';
	import { useSidebar } from '@/components/ui/sidebar';

	interface Props {
		content: Snippet<[{ isMobile: boolean }]>;
	}

	let { content }: Props = $props();

	const sidebar = useSidebar();
</script>

{@render content({ isMobile: sidebar.isMobile })}
