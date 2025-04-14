<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts" generics="Data extends GenericRecord">
	import type { GenericRecord } from '@/utils/types';
	import type { SuperForm, FormPath } from 'sveltekit-superforms';
	import { fieldProxy } from 'sveltekit-superforms';
	import * as Form from '@/components/ui/form';
	import FieldWrapper from './parts/fieldWrapper.svelte';
	import FileManager from '@/components/ui-custom/fileManager.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import type { FieldOptions } from './types';
	import type { Writable } from 'svelte/store';
	import type { ComponentProps } from 'svelte';
	import { createFilesValidator } from './fileField';

	//

	type Props = {
		form: SuperForm<Data>;
		name: FormPath<Data>;
		options?: Partial<FieldOptions & Omit<ComponentProps<typeof Input>, 'type' | 'value'>>;
	};

	const { form, name, options = {} }: Props = $props();

	const multiple = $derived(options.multiple ?? false);
	const valueProxy = $derived(fieldProxy(form, name) as Writable<File | File[]>);

	const validator = $derived(
		createFilesValidator(form as SuperForm<GenericRecord>, name, multiple)
	);
</script>

<Form.Field {form} {name}>
	<FieldWrapper field={name} {options}>
		{#snippet children({ props })}
			<FileManager bind:data={$valueProxy} {validator} {multiple}>
				{#snippet children({ addFiles })}
					<Input
						{...options}
						{...props}
						placeholder="Upload a file"
						type="file"
						class="hover:bg-primary/10 file:bg-secondary-foreground file:text-secondary p-0 py-1 pl-1 file:mr-4 file:h-full file:rounded-md file:px-4 hover:cursor-pointer file:hover:cursor-pointer"
						onchange={(e) => {
							const fileList = e.currentTarget.files;
							if (fileList) addFiles([...fileList]);
							e.currentTarget.value = '';
						}}
					/>
					<!-- e.currentTarget.value = '' is needed to clear the file input -->
				{/snippet}
			</FileManager>
		{/snippet}
	</FieldWrapper>
</Form.Field>
