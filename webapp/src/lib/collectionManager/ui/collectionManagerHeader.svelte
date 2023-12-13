<script lang="ts">
	import { Button, Heading, Modal, P } from 'flowbite-svelte';
	import { getRecordsManagerContext } from '../collectionManager.svelte';
	import { Trash, XMark } from 'svelte-heros-v2';
	import { CreateRecord } from './recordActions';
	import PortalWrapper from '$lib/components/portalWrapper.svelte';
	import type { ComponentProps } from 'svelte';

	//

	export let headingTag: ComponentProps<Heading>['tag'] = 'h4';
	export let showDeleteModal = false;
	export let description: string | null = null;
	export let hideCreateButton = false;

	const { collection, selectionManager, dataManager } = getRecordsManagerContext();
	const { recordService, loadRecords } = dataManager;
	const { selectedRecords, discardSelection } = selectionManager;

	async function deleteSelection() {
		if (!$selectedRecords.length) return;
		await Promise.all([...$selectedRecords.map((id) => recordService.delete(id))]);
		loadRecords();
		discardSelection();
	}
</script>

<div class="flex justify-between items-center mb-4">
	{#if $$slots.title}
		<slot name="title" />
	{:else}
		<Heading tag={headingTag}>{collection}</Heading>
	{/if}
	{#if description}
		<P class="text-slate-600 pt-6">{description}</P>
	{/if}
	<div class="shrink-0 flex space-x-4 items-center">
		{#if $selectedRecords.length > 0}
			<P><span class="font-bold">{$selectedRecords.length}</span> selected</P>
			<div class="flex space-x-2 items-center">
				<Button color="alternative" on:click={discardSelection}>
					<XMark size="20" />
					<span class="ml-1">Discard</span>
				</Button>
				<Button
					color="alternative"
					on:click={() => {
						showDeleteModal = true;
					}}
				>
					<Trash size="20" />
					<span class="ml-1">Delete</span>
				</Button>
			</div>
		{:else}
			<div class="flex space-x-2 items-center">
				{#if !hideCreateButton}
					<CreateRecord />
				{/if}
				<slot name="actions" />
			</div>
		{/if}
	</div>
</div>

<PortalWrapper>
	<Modal bind:open={showDeleteModal} title="Delete records" size="xs">
		<div class="text-center space-y-6">
			<P>
				Are you sure you want to delete <span class="font-bold">{$selectedRecords.length}</span> records?
			</P>
			<div class="flex gap-2 justify-center">
				<Button color="red" on:click={deleteSelection}>Delete</Button>
				<Button
					color="alternative"
					on:click={() => {
						showDeleteModal = false;
					}}>Cancel</Button
				>
			</div>
		</div>
	</Modal>
</PortalWrapper>
