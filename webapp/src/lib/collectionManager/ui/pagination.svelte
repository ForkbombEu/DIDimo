<script lang="ts">
	import { goto } from '$lib/i18n';
	import { Pagination } from 'flowbite-svelte';
	import { getRecordsManagerContext } from '../collectionManager.svelte';

	const { dataManager } = getRecordsManagerContext();
	const { perPage, currentPage, totalPages, totalItems } = dataManager;

	function handlePaginationClick(e: Event) {
		e.preventDefault();
		goto((e.target as any)?.href);
	}
	$: pages = Array.from({ length: $totalPages }, (_, i) => ({
		name: `${i + 1}`,
		href: `?page=${i + 1}`
	}));
</script>

{#if $totalPages > 1}
	<div class="my-5 flex flex-col items-center justify-center gap-2">
		<div class="text-sm text-gray-700 dark:text-gray-400">
			Showing <span class="font-semibold text-gray-900 dark:text-white"
				>{$perPage * Number($currentPage) - $perPage + 1}</span
			>
			to
			<span class="font-semibold text-gray-900 dark:text-white"
				>{Number($currentPage) == $totalPages ? $totalItems : $perPage * Number($currentPage)}</span
			>
			of <span class="font-semibold text-gray-900 dark:text-white">{$totalItems}</span> Entries
		</div>

		<div class="flex w-full justify-center">
			<Pagination
				{pages}
				activeClass="bg-blue-500 text-white"
				on:previous={(e) => {
					e.preventDefault();
					if (Number($currentPage) - 1 < 1) return;
					goto(`?page=${Number($currentPage) - 1}`);
				}}
				on:next={(e) => {
					e.preventDefault();
					if (Number($currentPage) + 1 > $totalPages) return;
					goto(`?page=${Number($currentPage) + 1}`);
				}}
				on:click={handlePaginationClick}
			/>
		</div>
	</div>
{/if}
