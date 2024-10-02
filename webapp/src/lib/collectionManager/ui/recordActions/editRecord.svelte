<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { createTypeProp } from '$lib/utils/typeProp';

	import type { PBResponse } from '$lib/utils/types';
	import { RecordForm, type FieldsSettings } from '$lib/recordForm';

	import PortalWrapper from '$lib/components/portalWrapper.svelte';
	import { Button, Modal } from 'flowbite-svelte';
	import { Pencil } from 'svelte-heros-v2';
	import { getRecordsManagerContext } from '../../collectionManager.svelte';

	//

	type RecordGeneric = $$Generic<PBResponse>;
	export let recordType = createTypeProp<RecordGeneric>();
	recordType;

	export let record: RecordGeneric;
	export let formSettings: Partial<FieldsSettings<RecordGeneric>> = {};

	//

	let { formFieldsSettings } = getRecordsManagerContext<RecordGeneric>();
	let { base, edit } = formFieldsSettings;

	let fieldsSettings = {
		...base,
		...edit,
		...formSettings
	};

	//

	let open = false;

	function openModal() {
		open = true;
	}
	function closeModal() {
		open = false;
	}
</script>

<slot {openModal}>
	<Button class="!p-2" color="alternative" on:click={openModal}>
		<Pencil size="20" />
	</Button>
</slot>

<PortalWrapper>
	<Modal bind:open title="Edit record" size="md">
		<div class="w-full">
			<RecordForm
				collection={record.collectionId}
				recordId={record.id}
				initialData={record}
				{fieldsSettings}
				on:success={closeModal}
			/>
		</div>
	</Modal>
</PortalWrapper>
