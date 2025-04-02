<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';

	let { data } = $props();
	const { news } = $derived(data);
</script>

<PageTop>
	<div class="flex flex-col gap-2">
		<T tag="h1">{news.title}</T>
		<T tag="p" class="text-card-foreground">{@html news.summary}</T>
		<T tag="small" class="text-muted-foreground">{new Date(news.updated).toLocaleString()}</T>
	</div>
	<div class="flex flex-row items-center justify-start gap-2">
		<a href={news.diff} target="_blank">
			<Button size="sm">{m.differences()}</Button>
		</a>
		<a href={news.refer} target="_blank">
			<Button size="sm">{m.referrer()}</Button>
		</a>
	</div>
</PageTop>

<PageContent class="grow bg-secondary" contentClass="flex gap-12 items-start">
	<div class="prose lg:prose-lg xl:prose-xl">
		{@html news.news}
	</div>
</PageContent>
