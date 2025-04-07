// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import type { CollectionFieldModeProp } from '@/collections-components/collectionField.svelte';
import type { CollectionName } from '@/pocketbase/collections-models';
import type {
	CollectionFormData,
	CollectionRecords,
	CollectionResponses,
	RecordIdString,
	CollectionRelatedCollections as Related
} from '@/pocketbase/types';
import type { PocketbaseQueryExpandOption } from '@/pocketbase/query';
import type { GenericRecord, KeyOf, MaybePromise } from '@/utils/types';
import type { Snippet } from 'svelte';
import type { FormPath, SuperForm } from 'sveltekit-superforms';
import type { CollectionInputRecordProps } from '../types';
import type { FormSnippets } from '@/forms/form.svelte';
import type { FormOptions } from '@/forms/form';

/* Props */

export type CollectionFormProps<C extends CollectionName> = CollectionFormOptions<C> &
	FormSnippets & {
		collection: C;
		recordId?: RecordIdString;
		initialData?: Partial<CollectionRecords[C]>;
	};

export type CollectionFormOptions<C extends CollectionName> = {
	onSuccess?: OnCollectionFormSuccess<C>;
	fieldsOptions?: Partial<FieldsOptions<C>>;
	uiOptions?: UIOptions;
	superformsOptions?: FormOptions<CollectionFormData[C]>;
};

/* On success */

export type CollectionFormMode = 'create' | 'edit';

type OnCollectionFormSuccess<C extends CollectionName> = (
	record: CollectionResponses[C],
	mode: CollectionFormMode
) => MaybePromise<void>;

/* Fields Options */

export type FieldsOptions<C extends CollectionName, R = CollectionFormData[C]> = {
	labels: { [K in keyof R]?: string };
	descriptions: { [K in keyof R]?: string };
	placeholders: { [K in keyof R]?: string };
	order: Array<KeyOf<R>>;
	exclude: Array<KeyOf<R>>;
	hide: { [K in keyof R]?: R[K] };
	defaults: { [K in keyof R]?: R[K] };
	relations: {
		[K in keyof Related[C]]?: RelationFieldOptions<Related[C][K] & CollectionName>;
	};
	snippets: { [K in keyof R]?: FieldSnippet<C> };
};

export type RelationFieldOptions<C extends CollectionName> = CollectionFieldModeProp &
	CollectionInputRecordProps<C, PocketbaseQueryExpandOption<C>>;

export type FieldSnippetOptions<C extends CollectionName, T = CollectionFormData[C]> = {
	form: SuperForm<T & GenericRecord>;
	field: FormPath<T & GenericRecord>;
};

export type FieldSnippet<C extends CollectionName, T = CollectionFormData[C]> = Snippet<
	[FieldSnippetOptions<C, T>]
>;

/* UI Options */

export type UIOptions = {
	hideRequiredIndicator?: boolean;
	showToastOnSuccess?: boolean;
	toastText?: string;
};
