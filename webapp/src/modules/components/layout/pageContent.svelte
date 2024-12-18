<script lang="ts">
	import { cn } from '@/components/ui/utils';

	interface Props {
		//
		layout?: 'vertical' | 'horizontal';
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { layout = 'vertical', class: className = '', children }: Props = $props();

	//

	let classes = $derived(
		cn('px-10 py-8', className, {
			'flex items-start gap-8': layout == 'horizontal',
			'space-y-8': layout == 'vertical'
		})
	);
</script>

<div class={classes}>
	{@render children?.()}
</div>
