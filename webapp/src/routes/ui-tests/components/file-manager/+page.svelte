<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import FileManager from '@/components/ui-custom/fileManager.svelte';
	import T from '@/components/ui-custom/t.svelte';
	import { createDummyFile } from '@/utils/other';

	let files: File[] = $state([
		createDummyFile({ filename: 'ciao' }),
		createDummyFile({ filename: 'np' }),
		createDummyFile({ filename: 'ok' })
	]);
</script>

<div class="space-y-8">
	<div class="space-y-4">
		<T tag="h3">Base file manager</T>
		<FileManager bind:data={files} />
	</div>
	<div class="space-y-4">
		<T tag="h3">Multiple file manager</T>
		<FileManager bind:data={files} multiple />
	</div>
</div>
