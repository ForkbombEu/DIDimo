<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { m } from '@/i18n';
	import T from '@/components/ui-custom/t.svelte';
	import Card from '@/components/ui-custom/card.svelte';
	import { Badge } from '@/components/ui/badge';
	import type { WalletsResponse } from '@/pocketbase/types';

	import { Separator } from '@/components/ui/separator';
	import A from '@/components/ui-custom/a.svelte';
	import { type ConformanceCheck } from '../../routes/my/services-and-products/wallet-form-checks-table.svelte';
	import { cn } from '@/components/ui/utils';
	import Avatar from '@/components/ui-custom/avatar.svelte';
	import { pb } from '@/pocketbase';

	type Props = {
		app: WalletsResponse;
		class?: string;
	};

	const { app, class: className = '' }: Props = $props();

	const logo = $derived(pb.files.getURL(app, app.logo));
</script>

<Card
	class={cn(
		'block overflow-auto rounded-lg border border-primary bg-card text-card-foreground shadow-sm ring-primary transition-all hover:-translate-y-2 hover:ring-2 ',
		className
	)}
>
	{@const conformanceChecks = app.conformance_checks as ConformanceCheck[]}
	<div class="space-y-4 overflow-scroll">
		<div class="flex flex-row justify-between">
			<div>
				<div class="flex items-center gap-2">
					{#if logo}
						<Avatar src={logo} class="!rounded-sm" hideIfLoadingError />
					{/if}
					<T class="font-bold">
						{#if !app.published}
							{app.name}
						{:else}
							<A href="/apps/{app.id}">{app.name}</A>
						{/if}
					</T>
				</div>
				<T class="mt-1 text-xs text-gray-400">
					{app.description}
				</T>
			</div>
		</div>

		<Separator />

		<div class="flex flex-wrap gap-2">
			{#if conformanceChecks.length > 0}
				{#each conformanceChecks as check}
					<Badge variant={check.status === 'success' ? 'secondary' : 'destructive'}>
						{check.test}
					</Badge>
				{/each}
			{:else}
				<T class="text-gray-300">{m.No_conformance_checks_available()}</T>
			{/if}
		</div>
	</div>
</Card>
