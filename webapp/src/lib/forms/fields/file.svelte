<script lang="ts" context="module">
	import type { LabelOption } from './types';
	import type { HTMLInputAttributes } from 'svelte/elements';

	export type FormInputOptions = Partial<HTMLInputAttributes> & LabelOption;

	//

	export function isFileArray(value: unknown): value is File[] {
		return Array.isArray(value) && value.every((item) => item instanceof File);
	}

	export function isFile(value: unknown): value is File {
		return Boolean(value) && value instanceof File;
	}
</script>

<script lang="ts">
	import type { z } from 'zod';
	import type { FormPath, ZodValidation } from 'sveltekit-superforms';
	import { formFieldProxy, type SuperForm } from 'sveltekit-superforms/client';

	import { Badge, Fileupload, Helper, Listgroup, ListgroupItem } from 'flowbite-svelte';
	import FieldError, { fieldHasErrors } from './fieldParts/fieldError.svelte';
	import FieldLabel from './fieldParts/fieldLabel.svelte';
	import ListgroupItemButton from '$lib/components/listgroupItemButton.svelte';

	type T = $$Generic<AnyZodObject>;

	export let superform: SuperForm<ZodValidation<T>, any>;
	export let field: FormPath<z.infer<T>>;
	export let options: FormInputOptions = {};

	let { multiple = false, accept = '', label = '' } = options;

	const { validate } = superform;
	const { value, errors } = formFieldProxy(superform, field as string);

	$: hasErrors = fieldHasErrors($errors);
	const oldFiles = ([] as File[]).concat($value);

	//

	function removeFile() {
		if (!isFile($value)) return;
		$value = undefined;
		validate(field as any);
	}

	function removeFileFromArray(file: File) {
		if (!isFileArray($value)) return;
		$value = $value.filter((item) => item !== file);
		validate(field as any);
	}

	//

	let rejectedFiles: { file: File; errors: string[] }[] = [];

	//

	async function handleFileSelect(e: Event) {
		const fileList = (e.target as HTMLInputElement)?.files;
		if (!fileList) return;

		let data = !multiple ? fileList[0] : ([...($value as File[]), ...fileList] as any);
		rejectedFiles = [];

		const errors = await validate(field as any, { value: data, update: false });
		if (!errors) {
			await validate(field as any, { value: data, taint: true });
			return;
		}

		if (!multiple && data instanceof File) {
			rejectedFiles = [{ file: data, errors }];
		}

		if (multiple && Array.isArray(data)) {
			const errorsRecord = errors as unknown as Record<string, string[]>;
			const rejectedFilesIndexes = Object.keys(errorsRecord).map((key) => Number(key));
			for (const index of rejectedFilesIndexes) {
				rejectedFiles = [
					...rejectedFiles,
					{
						file: data[index],
						errors: errorsRecord[index]
					}
				];
			}
			data = data.filter((_, index) => !rejectedFilesIndexes.includes(index));
			await validate(field as any, { value: data, taint: true });
		}
	}
</script>

<div class="space-y-3">
	{#if label}
		<FieldLabel {field} text={label} />
	{/if}

	<Fileupload
		{...options}
		id={field}
		name={field}
		{multiple}
		{accept}
		data-invalid={hasErrors}
		on:change={handleFileSelect}
	/>

	{#if isFile($value) || (isFileArray($value) && $value.length > 0)}
		<div>
			<Helper>FILES</Helper>

			<Listgroup>
				{#if !multiple && isFile($value)}
					{@const isNew = !oldFiles.includes($value)}
					<ListgroupItemButton on:remove={removeFile}>
						<div>
							{$value.name}
							{#if isNew}
								<Badge color="green">new</Badge>
							{/if}
						</div>
					</ListgroupItemButton>
				{/if}

				{#if multiple && isFileArray($value)}
					{#each $value as file}
						{@const isNew = !oldFiles.includes(file)}
						<ListgroupItemButton
							on:remove={() => {
								removeFileFromArray(file);
							}}
						>
							<div>
								{file.name}
								{#if isNew}
									<Badge color="green">new</Badge>
								{/if}
							</div>
						</ListgroupItemButton>
					{/each}
				{/if}
			</Listgroup>
		</div>
	{/if}

	<FieldError {field} />

	{#if rejectedFiles.length > 0}
		<div class="space-y-1">
			<Helper>REJECTED FILES</Helper>
			<Listgroup>
				{#each rejectedFiles as { file, errors }}
					<ListgroupItem>
						<div>
							{file.name}
						</div>
						<div>
							{#each errors as error}
								<Helper color="red">{error}</Helper>
							{/each}
						</div>
					</ListgroupItem>
				{/each}
			</Listgroup>
		</div>
	{/if}
</div>
