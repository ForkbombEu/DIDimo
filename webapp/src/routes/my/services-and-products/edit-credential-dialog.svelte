<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { CollectionForm } from '@/collections-components';
	import * as Dialog from '@/components/ui/dialog';
	import { buttonVariants } from '@/components/ui/button';
	import { Pencil } from 'lucide-svelte';
	import type { CredentialsRecord } from '@/pocketbase/types';

	type Props = {
		credential: CredentialsRecord;
	};

	let { credential }: Props = $props();

	let open = $state(false);
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class={buttonVariants({ variant: 'default' })}>
		<Pencil />
	</Dialog.Trigger>

	<Dialog.Content class=" sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Credential {credential.key}</Dialog.Title>
		</Dialog.Header>

		<div class="pt-8">
			<CollectionForm
				collection="credentials"
				recordId={credential.id}
				initialData={credential}
				fieldsOptions={{
					exclude: [
						'format',
						'issuer_name',
						'type',
						'name',
						'locale',
						'logo',
						'description',
						'credential_issuer',
						'json',
						'key'
					],
					order: ['published']
				}}
				onSuccess={() => {
					open = false;
				}}
			/>
		</div>
	</Dialog.Content>
</Dialog.Root>
