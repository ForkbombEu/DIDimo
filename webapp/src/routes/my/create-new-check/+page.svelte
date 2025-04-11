<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { CollectionForm } from '@/collections-components';
	import CollectionManager from '@/collections-components/manager/collectionManager.svelte';
</script>

<CollectionManager collection="custom_checks">
	{#snippet top({ Header })}
		<Header />
	{/snippet}

	{#snippet records({ records })}
		{#each records as record}
			<pre class="bg-white p-4">{JSON.stringify(record, null, 4)}</pre>
		{/each}
	{/snippet}
</CollectionManager>
