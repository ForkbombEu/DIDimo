<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	import type { HTMLInputAttributes } from 'svelte/elements';

	export type FormHiddenOptions = Partial<HTMLInputAttributes>;
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPathLeaves<z.infer<T>>;
	export let options: FormHiddenOptions = {};

	const { value } = formFieldProxy(superform, field as string);
</script>

<input {...options} type="hidden" id={field} name={field} bind:value={$value} />
