<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	/**
	 * Superforms validation errors come in different shapes:
	 * • string[] (the most common)
	 * • {"1": string[], "2": string[], ...} (when using arrays)
	 * • {"key": string[], ...} (when using nested fields)
	 *
	 * Here we try to detect the shape of the error data.
	 */

	export function isNonNullable(value: unknown): value is NonNullable<unknown> {
		return value !== undefined && value !== null;
	}

	export function isBaseError(errorData: unknown): errorData is string[] {
		return Array.isArray(errorData) && errorData.length > 0;
	}

	export function isNestedError(errorData: unknown): errorData is Record<string, string[]> {
		return (
			isNonNullable(errorData) &&
			!isBaseError(errorData) &&
			typeof errorData === 'object' &&
			Object.values(errorData).some((value) => isBaseError(value))
		);
	}

	export function fieldHasErrors(errorData: unknown): boolean {
		return isBaseError(errorData) || isNestedError(errorData);
	}
</script>

<script lang="ts">
	import { formFieldProxy } from 'sveltekit-superforms/client';
	import { getFormContext } from '../../form.svelte';
	import { Helper } from 'flowbite-svelte';

	export let field: string;

	const { superform } = getFormContext();
	const { errors } = formFieldProxy(superform, field);
</script>

{#if isBaseError($errors)}
	<div class="space-y-1">
		{#each $errors as error}
			<Helper color="red">{error}</Helper>
		{/each}
	</div>
{/if}

{#if isNestedError($errors)}
	<div class="space-y-2">
		{#each Object.entries($errors) as [key, errors]}
			{#if isBaseError(errors)}
				<div class="space-y-1">
					{#if key !== '_errors'}
						<Helper color="red"><span class="font-bold">{key}</span></Helper>
					{/if}
					{#each errors as error}
						<Helper color="red">{error}</Helper>
					{/each}
				</div>
			{/if}
		{/each}
	</div>
{/if}
