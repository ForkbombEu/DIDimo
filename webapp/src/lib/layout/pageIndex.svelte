<script lang="ts" module>
	import type { IconComponent } from '@/components/types';

	export interface IndexItem {
		icon: IconComponent;
		anchor: string;
		label: string;
	}
</script>

<script lang="ts">
	interface Props {
		sections: IndexItem[];
	}

	let { sections }: Props = $props();
</script>

<ul class="space-y-4">
	{#each sections as section}
		<li>
			<a href={`#${section.anchor}`} class="flex items-center gap-2 hover:underline">
				<section.icon class="size-4" />
				{section.label}
			</a>
		</li>
	{/each}
</ul>
