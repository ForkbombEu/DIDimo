<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" context="module">
	import { getContext } from 'svelte';
	import { normalizeError, type ClientResponseErrorData } from '$lib/errorHandling';
	import type { AnyZodObject, ZodEffects } from 'zod';
	import type { UnwrapEffects, ZodValidation } from 'sveltekit-superforms';
	import {
		superForm,
		setMessage,
		setError,
		type FormOptions,
		type SuperForm,
		superValidateSync
	} from 'sveltekit-superforms/client';
	import type z from 'zod';

	//

	export const FORM_KEY = Symbol('form');

	export type FormContext<T extends AnyZodObject> = {
		superform: SuperForm<ZodValidation<T>, ClientResponseErrorData>;
		showRequiredIndicator: boolean;
	};

	export function getFormContext<T extends AnyZodObject>(): FormContext<T> {
		return getContext(FORM_KEY);
	}

	//

	export type FormSettings<T extends AnyZodObject = AnyZodObject> = FormOptions<
		ZodValidation<T>,
		unknown
	>;

	export type SubmitFunction<T extends AnyZodObject> = NonNullable<FormSettings<T>['onUpdate']>;

	export function createForm<T extends AnyZodObject>(
		schema: T | ZodEffects<T>,
		submitFunction: SubmitFunction<T> = async () => {},
		initialData: Partial<z.infer<T>> | undefined = undefined,
		options: FormSettings<T> = {}
	) {
		const form = superValidateSync(initialData, schema, { errors: false });
		return superForm<ZodValidation<T>, ClientResponseErrorData>(form, {
			SPA: true,
			applyAction: false,
			scrollToError: 'smooth',
			// @ts-ignore
			validators: schema,
			dataType: 'json',
			onUpdate: async (input) => {
				try {
					if (input.form.valid) await submitFunction(input);
				} catch (e) {
					let error = normalizeError(e);
					for (const [key, value] of Object.entries(error.data)) {
						if (Boolean(input.form.data[key])) {
							setError(input.form, key as any, value.message);
							delete error.data[key];
						}
					}
					setMessage(input.form, error);
				}
			},
			taintedMessage: null,
			...options
		});
	}

	type SuperFormAllErrors = {
		path: string;
		messages: string[];
	}[];

	export function formHasErrors(allErrors: SuperFormAllErrors): boolean {
		const checks: boolean[] = [];
		if (allErrors.length === 0) return false;
		else {
			for (const error of allErrors) {
				if (error.messages.length > 0) checks.push(true);
				// if (typeof value === 'object') checks.push(formHasErrors(value));
				// else if (Boolean(value)) checks.push(true);
			}
		}
		return checks.some((check) => check);
	}

	//

	export function createFormData(data: Record<string, unknown>) {
		const formData = new FormData();
		for (const [key, value] of Object.entries(data)) {
			if (value === null || value === undefined) {
				// Needed otherwise pb complains about "bad formatting", especially for null files
				formData.append(key, '');
			} else if (Array.isArray(value)) {
				// Special case for empty arrays, cause they can't be represented in formData
				if (value.length === 0) {
					formData.append(key, '');
				} else {
					for (const item of value) {
						formData.append(key, item);
					}
				}
			} else {
				formData.append(key, value as string | File);
			}
		}
		return formData;
	}
</script>

<script lang="ts">
	import { setContext } from 'svelte';
	import { Spinner, Modal } from 'flowbite-svelte';
	import PortalWrapper from '$lib/components/portalWrapper.svelte';

	type T = $$Generic<AnyZodObject>;

	//

	export let superform: SuperForm<UnwrapEffects<T>, any>;

	export let showRequiredIndicator = false;
	export let className = 'space-y-8';

	let enctype = superform.options.dataType == 'form' ? 'multipart/form-data' : undefined;

	//

	const { enhance, delayed } = superform;
	setContext<FormContext<T>>(FORM_KEY, { superform, showRequiredIndicator });
</script>

<form class={className} method="post" use:enhance {enctype}>
	<slot />
</form>

{#if $delayed}
	<PortalWrapper>
		<Modal open={$delayed} dismissable={false}>
			<Spinner />
		</Modal>
	</PortalWrapper>
{/if}
