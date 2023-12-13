<script lang="ts" context="module">
	export type ColumnType =
        | 'avatar'
		| 'image'
		| 'short text'
		| 'long text'
		| 'actions'
		| 'date'
		| 'file'
		| 'checkbox';
</script>

<script lang="ts">
	export let columns: ColumnType[] = [
        'avatar',
		'image',
		'short text',
		'long text',
		'actions',
		'date',
		'file',
		'checkbox'
	];
	export let rows: number = 5;

	function getColumnClass(columnType: ColumnType) {
		switch (columnType) {
			case 'avatar':
				return 'w-16 h-16 bg-gray-200 rounded-full animate-pulse duration-700';
            case 'image':
                return 'w-16 h-16 bg-gray-200 animate-pulse duration-700';
			case 'short text':
				return 'w-24 h-4 bg-gray-200 animate-pulse rounded-lg duration-700';
			case 'long text':
				return 'h-2 bg-gray-200 animate-pulse rounded-lg duration-700';
			case 'actions':
				return 'w-4 h-4 bg-gray-200 animate-pulse rounded-lg duration-700';
			case 'date':
				return 'w-16 h-4 bg-gray-200 animate-pulse rounded-lg duration-700';
			case 'file':
				return 'w-24 h-4 bg-gray-200 animate-pulse rounded-lg duration-700';
			case 'checkbox':
				return 'w-4 h-4 bg-gray-200 animate-pulse duration-700';
			default:
				return 'w-24 h-4 bg-gray-200 animate-pulse rounded-lg duration-700';
		}
	}
</script>

<div class="p-4 space-y-1 rounded-lg">
	<div
		class="flex p-4 justify-between items-center bg-gray-400 animate-pulse duration-500 rounded-lg"
	>
		{#each Array(columns.length) as _, i}
			<div class="w-32">
				{#if columns[i] === 'checkbox'}
					<div class="w-4 h-4 bg-gray-200 animate-pulse duration-700" />
				{:else if columns[i] === 'actions'}
					<div class="w-24" />
				{:else}
					<div class="w-24 h-4 bg-gray-200 animate-pulse duration-700 rounded-lg" />
				{/if}
			</div>
		{/each}
	</div>
	{#each Array(rows) as _}
		<div
			class="flex p-4 justify-between items-center bg-gray-400 animate-pulse duration-500 rounded-lg"
		>
			{#each columns as columnType}
				<div class="w-32">
					{#if columnType === 'long text'}
						<div class="flex flex-col space-y-1">
							{#each Array(5) as _, i}
								<div class={`${getColumnClass(columnType)} w-${(i + 1) % 2 !== 0 ? '32' : '24'}`} />
							{/each}
						</div>
					{:else if columnType === 'actions'}
						<div class="flex space-x-2">
							{#each Array(3) as _, i}
								<div class={getColumnClass(columnType)} />
							{/each}
						</div>
					{:else}
						<div class={getColumnClass(columnType)} />
					{/if}
				</div>
			{/each}
		</div>
	{/each}
</div>

<!--
    @component TableSkeleton
    @type ColumnType - 'avatar' | 'image' | 'short text' | 'long text' | 'actions' | 'date' | 'file' | 'checkbox'
    @description Skeleton for list table with columns and rows
    @prop columns - array of column types, default: ['avatar','image', 'short text', 'long text', 'actions', 'date', 'file']
    @prop rows - number of rows, default: 5
-->
