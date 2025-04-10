<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import LoadingDialog from '@/components/ui-custom/loadingDialog.svelte';
	import { pb } from '@/pocketbase';
	import type { BaseSystemFields } from '@/pocketbase/types';
	import type { Snippet } from 'svelte';

	type Props = {
		record: Omit<BaseSystemFields, 'expand'> & { published: boolean };
		button: Snippet<[{ togglePublish: () => void; label: string }]>;
		onSuccess?: () => void;
	};

	let { record, button, onSuccess }: Props = $props();

	let loading = $state(false);

	async function publish(value: boolean) {
		loading = true;
		await pb.collection(record.collectionName).update(record.id, {
			published: value
		});
		onSuccess?.();
		loading = false;
	}

	const togglePublish = $derived(() => {
		publish(!record.published);
	});

	const label = $derived.by(() => {
		return record.published ? 'Unpublish' : 'Publish';
	});
</script>

{@render button({ togglePublish, label })}

{#if loading}
	<LoadingDialog />
{/if}
