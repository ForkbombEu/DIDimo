<script lang="ts">
	import type { PBResponse } from '$lib/utils/types';
	import { getRecordsManagerContext } from '../../collectionManager.svelte';

	import { Alert, Button, Modal, P } from 'flowbite-svelte';
	import { Trash, XMark } from 'svelte-heros-v2';
	import PortalWrapper from '$lib/components/portalWrapper.svelte';

	type RecordGeneric = $$Generic<PBResponse>;
	export let record: RecordGeneric;

	const { dataManager } = getRecordsManagerContext();
	const { loadRecords, recordService } = dataManager;

	async function deleteRecord() {
		try {
			await recordService.delete(record.id);
			loadRecords();
			open = false;
		} catch (e) {
			error = e;
			if (e && typeof e == 'object' && 'response' in e) error = e.response;
		}
	}

	let open = false;
	let error: unknown;
</script>

<Button
	class="!p-2"
	color="alternative"
	on:click={() => {
		open = true;
	}}
>
	<Trash size="20" />
</Button>

<PortalWrapper>
	<Modal
		bind:open
		title="Delete record"
		size="xs"
		on:close={() => {
			error = undefined;
		}}
	>
		<div class="text-center space-y-6">
			<P>Are you sure you want to delete this record?</P>
			{#if error}
				<Alert color="red" dismissable>
					<P weight="bold" color="text-red-800 dark:text-red-400">Error</P>
					<pre class="text-left">{JSON.stringify(error, null, 2)}</pre>
				</Alert>
			{/if}
			<div class="flex gap-2 justify-center">
				<Button color="red" on:click={deleteRecord}>
					<Trash size="20" /><span class="ml-1">Delete</span>
				</Button>
				<Button
					color="alternative"
					on:click={() => {
						open = false;
					}}
				>
					<XMark size="20" />
					<span class="ml-1">Cancel</span>
				</Button>
			</div>
		</div>
	</Modal>
</PortalWrapper>
