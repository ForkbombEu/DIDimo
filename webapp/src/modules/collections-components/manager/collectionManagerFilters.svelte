<script lang="ts">
	import { buttonVariants } from '@/components/ui/button';
	import {
		getCollectionManagerContext,
		type Filter,
		FilterSchema,
		FilterGroupSchema
	} from './collectionManagerContext';
	import * as Popover from '@/components/ui/popover/index.js';
	import * as Sheet from '@/components/ui/sheet/index.js';

	//

	type Props = {
		modalType?: 'popover' | 'sheet';
	};

	let { modalType = 'popover' }: Props = $props();

	const { filters, manager } = getCollectionManagerContext();

	// TODO - Use manager.query as source for the selected filters
	let selectedFilters = $state<Filter[]>([]);

	$effect(() => {
		manager.query.setFilters(selectedFilters.map((f) => f.expression));
	});
</script>

{#if modalType === 'popover'}
	<Popover.Root>
		<Popover.Trigger class={buttonVariants({ variant: 'outline' })}>Open</Popover.Trigger>
		<Popover.Content class="w-80">
			{@render content()}
		</Popover.Content>
	</Popover.Root>
{/if}

{#if modalType === 'sheet'}
	<Sheet.Root>
		<Sheet.Trigger>Open</Sheet.Trigger>
		<Sheet.Content>
			<Sheet.Header>
				<Sheet.Title>Are you sure absolutely sure?</Sheet.Title>
				<Sheet.Description>
					This action cannot be undone. This will permanently delete your account and
					remove your data from our servers.
				</Sheet.Description>
			</Sheet.Header>
			{@render content()}
		</Sheet.Content>
	</Sheet.Root>
{/if}

{#snippet content()}
	<ul>
		{#each filters as filterData}
			{@const filter = FilterSchema.safeParse(filterData)}
			{@const filterGroup = FilterGroupSchema.safeParse(filterData)}

			{#if filter.success}
				<li>
					{@render filterInput(filter.data)}
				</li>
			{:else if filterGroup.success}
				<ul>
					<li>{filterGroup.data.name}</li>
					{#each filterGroup.data.filters as filter}
						<li>
							{@render filterInput(filter)}
						</li>
					{/each}
				</ul>
			{/if}
		{/each}
	</ul>
{/snippet}

{#snippet filterInput(f: Filter)}
	<label>
		<input type="checkbox" name={f.id} id={f.id} value={f} bind:group={selectedFilters} />
		{f.name}
	</label>
{/snippet}
