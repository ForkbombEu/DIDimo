<script lang="ts" module>
	import * as Avatar from '@/components/ui/avatar';
	import type { ComponentProps } from 'svelte';

	export type AvatarProps = ComponentProps<typeof Avatar.Root> & {
		src?: string;
		alt?: string;
		fallback?: string;
	};
</script>

<script lang="ts">
	import { m } from '@/i18n';
	import { cn } from '../ui/utils';

	const { src, alt, fallback, ...rest }: AvatarProps = $props();
</script>

<Avatar.Root {...rest} class={cn(rest.class, 'overflow-hidden')}>
	{#if src}
		<Avatar.Image {src} alt={alt ?? m.Avatar()} />
	{/if}
	{#if fallback}
		<Avatar.Fallback class="rounded-none text-[80%]">{fallback}</Avatar.Fallback>
	{/if}
</Avatar.Root>
