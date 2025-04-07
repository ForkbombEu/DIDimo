<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" module>
	import type { GenericRecord } from '@/utils/types';
</script>

<script lang="ts" generics="T extends GenericRecord">
	import SuperDebug, { type SuperForm } from 'sveltekit-superforms';

	interface Props {
		form: SuperForm<T>;
	}

	let { form }: Props = $props();
	const { form: formData } = form;
</script>

<SuperDebug data={$formData}></SuperDebug>
