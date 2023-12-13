<script lang="ts">
	import { ListgroupItem } from 'flowbite-svelte';
	import IconButton from './iconButton.svelte';
	import { Eye } from 'svelte-heros-v2';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	const onRemove = () => {
		dispatch('remove');
	};

	const onShow = () => {
		dispatch('show');
	};
</script>

<ListgroupItem>
	<div class="flex justify-between items-center">
		<slot />
		<div>
			<IconButton icon={Eye} on:click={onShow} />
			<IconButton on:click={onRemove} />
		</div>
	</div>
</ListgroupItem>
