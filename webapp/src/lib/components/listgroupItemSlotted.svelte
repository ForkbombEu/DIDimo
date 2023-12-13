<script lang="ts">
	import { ListgroupItem } from 'flowbite-svelte';
</script>

<ListgroupItem class="!py-1 !pr-1 !pl-3">
	<div class="flex justify-between items-center flex-nowrap gap-4">
		<div class="grow w-0 overflow-hidden whitespace-nowrap">
			<slot name="left" />
		</div>
		<div class="flex justify-end items-center gap-1">
			<slot name="right" />
		</div>
	</div>
</ListgroupItem>
