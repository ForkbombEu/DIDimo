<script lang="ts">
	import type { PBResponse } from '$lib/utils/types';
	import { createTypeProp } from '$lib/utils/typeProp';
	import { Collections, type AuthorizationsResponse } from '$lib/pocketbase/types';
	import { RecordForm } from '$lib/recordForm';
	import { Button, Modal, Spinner, P } from 'flowbite-svelte';
	import { currentUser, pb } from '$lib/pocketbase';
	import { ArrowLeft, Share, Trash } from 'svelte-heros-v2';
	import { createEventDispatcher } from 'svelte';

	type RecordGeneric = $$Generic<PBResponse>;
	export let record: RecordGeneric;

	const dispatch = createEventDispatcher<{ add: undefined; remove: undefined }>();

	/* */

	let open = false;
	const openModal = () => {
		open = true;
	};

	/* Load */

	let authorization: AuthorizationsResponse | undefined;
	const authorizationRequest = loadAuthorization();
	const recordType = createTypeProp<AuthorizationsResponse>();

	async function loadAuthorization() {
		try {
			authorization = await pb
				.collection(Collections.Authorizations)
				.getFirstListItem<AuthorizationsResponse>(`record_id="${record.id}"`);
		} catch (e) {
			authorization = undefined;
		}
	}

	/* Success */

	function handleSuccess() {
		open = false;
		dispatch('add');
	}

	/* Remove */

	let removeAccess = false;
	let removeLoading = false;

	function toggleRemoveAccess() {
		removeAccess = !removeAccess;
	}

	async function removeAuthorization() {
		if (!authorization) return;
		removeLoading = true;
		await pb.collection(Collections.Authorizations).delete(authorization.id);
		open = false;
		dispatch('remove');
	}
</script>

<slot {openModal}>
	<Button
		class="!p-2"
		color="alternative"
		on:click={() => {
			openModal();
		}}
	>
		<Share size="20" />
	</Button>
</slot>

{#await authorizationRequest}
	<Spinner />
{:then response}
	<div class="fixed z-50">
		<Modal bind:open size="md" title="Share signature">
			<div class="w-full relative">
				{#if !removeAccess}
					<RecordForm
						{recordType}
						on:success={handleSuccess}
						initialData={authorization}
						recordId={authorization?.id}
						submitButtonText="Share"
						collection={Collections.Authorizations}
						fieldsSettings={{
							hide: {
								record_id: record.id,
								collection_id: record.collectionId,
								owner: $currentUser?.id
							},
							relations: {
								users: {
									displayFields: ['email']
								}
							}
						}}
					/>
					{#if authorization}
						<div class="absolute left-0 bottom-0">
							<Button color="red" outline on:click={toggleRemoveAccess}>
								<Trash size="20" />
								<span class="ml-2"> Remove access </span>
							</Button>
						</div>
					{/if}
				{:else if removeLoading}
					<Spinner />
				{:else}
					<P>Are you sure you want to remove all access to the signature?</P>
					<div class="flex justify-between mt-4">
						<Button class="space-x-2" color="alternative" on:click={toggleRemoveAccess}>
							<ArrowLeft size="20" />
							<span class="ml-2"> Undo </span>
						</Button>
						<Button class="space-x-2" color="red" on:click={removeAuthorization}>
							<Trash size="20" />
							<span class="ml-2"> Yes, remove access </span>
						</Button>
					</div>
				{/if}
			</div>
		</Modal>
	</div>
{/await}
