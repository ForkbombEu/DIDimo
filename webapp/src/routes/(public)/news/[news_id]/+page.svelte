<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import PageContent from '$lib/layout/pageContent.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { Button } from '@/components/ui/button';
	import { m } from '@/i18n';
	import HTML from '@/components/ui-custom/renderHTML.svelte';

	let { data } = $props();
	const { news } = $derived(data);
</script>

<PageTop>
	<div class="flex flex-col gap-4">
		<T tag="h1">{news.title}</T>
		<T tag="small" class="text-muted-foreground">{new Date(news.updated).toLocaleString()}</T>
		<HTML class="prose prose-sm xl:prose-lg" content={news.summary} />
	</div>
	<div class="flex flex-row items-center justify-start gap-2">
		{#if news.diff}
			<a href={news.diff} target="_blank">
				<Button size="sm">{m.differences()}</Button>
			</a>
		{/if}
		{#if news.refer}
			<a href={news.refer} target="_blank">
				<Button size="sm">{m.referrer()}</Button>
			</a>
		{/if}
	</div>
</PageTop>

<PageContent class="grow bg-secondary" contentClass="flex gap-12 items-start">
	<div class="prose prose-base lg:prose-lg xl:prose-xl">
		<HTML content={news.news} />
	</div>
</PageContent>
