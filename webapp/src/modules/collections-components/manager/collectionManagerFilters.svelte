<script lang="ts">
	import {
		getCollectionManagerContext,
		type Filter,
		FilterSchema,
		FilterGroupSchema
	} from './collectionManagerContext';

	const { filters, manager } = getCollectionManagerContext();

	let selectedFilters = $state<Filter[]>([]);

	$effect(() => {
		manager.query.setFilters(selectedFilters.map((f) => f.expression));
	});
</script>

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

{#snippet filterInput(f: Filter)}
	<label>
		<input type="checkbox" name={f.id} id={f.id} value={f} bind:group={selectedFilters} />
		{f.name}
	</label>
{/snippet}
