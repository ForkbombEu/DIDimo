<script lang="ts" module>
	import type { IconComponent } from '@/components/types';
	import A from '@/components/ui-custom/a.svelte';
	export interface IndexItem {
		icon: IconComponent;
		anchor: string;
		label: string;
	}
</script>

<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	type Props = {
		sections: IndexItem[];
	} & HTMLAttributes<HTMLUListElement>;

	let { sections, class: className, ...props }: Props = $props();
</script>

<ul class="space-y-4 {className}" {...props}>
	{#each sections as section}
		<li>
			<A href={`#${section.anchor}`} class="flex items-center gap-2 hover:underline">
				<section.icon class="size-4" />
				{section.label}
			</A>
		</li>
	{/each}
</ul>
