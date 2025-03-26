<script lang="ts">
	import type { SnippetFunction } from '@/components/types';
	import Icon from '@/components/ui-custom/icon.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import Button from '@/components/ui-custom/button.svelte'';
	import { ArrowUpRight } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		id: string;
		class?: string;
		title?: string;
		children?: Snippet;
		right?: Snippet<[{ Link: SnippetFunction<LinkProps> }]>;
	}

	interface LinkProps {
		href?: string;
		label: string;
		onclick?: () => void;
	}

	let { title, id, class: className = '', children, right }: Props = $props();
</script>

<div
	{id}
	class="mb-6 flex scroll-mt-5 items-center justify-between border-b border-secondary-foreground {className}"
>
	{#if title}
		<T tag="h2">{title}:</T>
	{/if}

	{@render children?.()}

	{@render right?.({ Link })}
</div>

{#snippet Link(props: LinkProps)}
	{@const { href, label, onclick } = props}
	<Button variant="link" {href} {onclick} class="gap-1 underline hover:no-underline">
		<T>{label}</T>
		<Icon src={ArrowUpRight} />
	</Button>
{/snippet}
