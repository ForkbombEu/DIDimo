<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { cn } from '@/components/ui/utils';
	import { m } from '@/i18n';
	import type { NewsResponse } from '@/pocketbase/types';

	type Props = {
		news: NewsResponse;
		class?: string;
	};

	const { news, class: className = '' }: Props = $props();
</script>

<div class={cn('flex flex-col gap-6 text-card-foreground', className)}>
	<div class="flex flex-col">
		<T tag="h3" class="block">{news.title}</T>
		<T tag="small" class="block">{@html news.summary}</T>
	</div>
	<T tag="small" class="block font-normal leading-snug">{@html `${news.news.slice(0, 200)}...`}</T
	>
	<a href="/news/{news.id}">
		<Button variant="outline" size="sm">
			{m.view_more()}
		</Button>
	</a>
</div>
