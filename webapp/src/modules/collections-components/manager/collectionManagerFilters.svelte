<script lang="ts">
	import { Button, buttonVariants, type ButtonVariant } from '@/components/ui/button';
	import {
		getCollectionManagerContext,
		type Filter,
		FilterSchema,
		FilterGroupSchema
	} from './collectionManagerContext';
	import * as Popover from '@/components/ui/popover/index.js';
	import * as Sheet from '@/components/ui/sheet/index.js';
	import { Checkbox } from '@/components/ui/checkbox';
	import { Label } from '@/components/ui/label';
	import T from '@/components/ui-custom/t.svelte';
	import type { Snippet } from 'svelte';
	import type { GenericRecord } from '@/utils/types';

	//

	type Props = {
		modalType?: 'popover' | 'sheet';
		triggerVariant?: ButtonVariant;
		children: Snippet;
		beforeFilters?: Snippet;
		afterFilters?: Snippet;
		trigger?: Snippet<[{ props: GenericRecord }]>;
	};

	let {
		children,
		modalType = 'popover',
		triggerVariant = 'outline',
		beforeFilters,
		afterFilters,
		trigger
	}: Props = $props();

	const { filters, manager } = getCollectionManagerContext();
</script>

{#if modalType === 'popover'}
	<Popover.Root>
		<Popover.Trigger class={buttonVariants({ variant: triggerVariant })}>
			{#snippet child({ props })}
				{@render childSnippet(props)}
			{/snippet}
		</Popover.Trigger>

		<Popover.Content class="w-80">
			{@render content()}
		</Popover.Content>
	</Popover.Root>
{/if}

{#if modalType === 'sheet'}
	<Sheet.Root>
		<Sheet.Trigger class={buttonVariants({ variant: triggerVariant })}>
			{#snippet child({ props })}
				{@render childSnippet(props)}
			{/snippet}
		</Sheet.Trigger>

		<Sheet.Content>
			{@render content()}
		</Sheet.Content>
	</Sheet.Root>
{/if}

{#snippet content()}
	{@render beforeFilters?.()}
	<ul class="space-y-4">
		{#each filters as filterData}
			{@const filter = FilterSchema.safeParse(filterData)}
			{@const filterGroup = FilterGroupSchema.safeParse(filterData)}

			{#if filter.success}
				<li>
					{@render filterInput(filter.data)}
				</li>
			{:else if filterGroup.success}
				<ul class="space-y-2">
					<li><T>{filterGroup.data.name}</T></li>
					{#each filterGroup.data.filters as filter}
						<li>
							{@render filterInput(filter)}
						</li>
					{/each}
				</ul>
			{/if}
		{/each}
	</ul>
	{@render afterFilters?.()}
{/snippet}

{#snippet filterInput(f: Filter)}
	<div class="flex items-center gap-2">
		<Checkbox
			name={f.id}
			id={f.id}
			value={f.expression}
			checked={manager.query.hasFilter(f.expression)}
			onCheckedChange={(v) => {
				if (v) {
					manager.query.addFilter(f.expression);
				} else {
					manager.query.removeFilter(f.expression);
				}
			}}
		/>
		<Label for={f.id} class="text-md">{f.name}</Label>
	</div>
{/snippet}

{#snippet childSnippet(props: GenericRecord)}
	{#if trigger}
		{@render trigger({ props })}
	{:else}
		<Button {...props} variant={triggerVariant}>
			{@render children()}
		</Button>
	{/if}
{/snippet}
