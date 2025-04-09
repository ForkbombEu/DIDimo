<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<!-- 
	! IMPORTANT !
	This page has `@` inside its name in order to "reset" the layout stack.
	Styles imported from the `/temporal.css` file conflict with the global styles.
-->

<script>
	import BackButton from '$lib/layout/back-button.svelte';
	import PageTop from '$lib/layout/pageTop.svelte';
	import T from '@/components/ui-custom/t.svelte';

	import TemporalI18nProvider from './components/temporal-i18n-provider.svelte';
	import TemporalWorkflow from './components/temporal-workflow.svelte';

	//

	let { data } = $props();
	const { workflowId, workflowResponse, eventHistory } = data;
</script>

<!--  -->

<!--
	! IMPORTANT !
	temporal.css is exported from @forkbombeu/temporal-ui
	but must be placed in static assets to work properly
-->
<svelte:head>
	<link rel="stylesheet" href="/temporal.css" />
</svelte:head>

<!--  -->

<div class="bg-[#220e7c] text-white">
	<div class="mx-auto max-w-screen-xl !px-2 md:!px-4 lg:!px-8">
		<BackButton href="/my/workflows" class="text-white">Back to workflows</BackButton>
	</div>
</div>

<PageTop contentClass="!space-y-0 !px-2 md:!px-4 lg:!px-8">
	<T tag="h1">Workflow {workflowId}</T>
</PageTop>

<div class="mx-auto max-w-screen-xl">
	<TemporalI18nProvider>
		<TemporalWorkflow {workflowResponse} {eventHistory} />
	</TemporalI18nProvider>
</div>
