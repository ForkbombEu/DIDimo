// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import type { GenericRecord } from '@/utils/types';
import {
	createInitialDataFromFields,
	createTestConfigFormSchema,
	type ConfigField,
	type ConfigFieldSpecific
} from './logic';
import type { SuperForm } from 'sveltekit-superforms';
import { createForm } from '@/forms';
import { zod } from 'sveltekit-superforms/adapters';
import { nanoid } from 'nanoid';
import { watch } from 'runed';
import { Record as R, pipe } from 'effect';
import { fromStore } from 'svelte/store';

//

type ConfigFields = ConfigField[] | ConfigFieldSpecific[];

type Getter<T> = () => T;

//

export class TestConfigFieldsForm {
	readonly form: SuperForm<GenericRecord>;
	readonly specificFields: ConfigFieldSpecific[];

	private overriddenFieldsIds = $state<string[]>([]);
	readonly currentOverriddenFields = $derived.by(() => {
		return this.fields
			.filter((f) => this.overriddenFieldsIds.includes(f.CredimiID))
			.map((f) => f as ConfigFieldSpecific);
	});

	readonly currentSharedFields = $derived.by(() => {
		return this.fields
			.filter((f) => this.sharedFieldsIds.includes(f.CredimiID))
			.filter((f) => !this.overriddenFieldsIds.includes(f.CredimiID));
	});

	isValid = $state(false);
	readonly currentData: GenericRecord;

	constructor(
		readonly id: string,
		readonly fields: ConfigFields,
		readonly sharedFieldsIds: string[],
		readonly sharedData: Getter<GenericRecord>
	) {
		this.fields = fields;

		this.specificFields = fields
			.filter((f) => !sharedFieldsIds.includes(f.CredimiID))
			.map((f) => f as ConfigFieldSpecific);

		//

		this.form = createForm({
			adapter: zod(createTestConfigFormSchema(fields)),
			initialData: createInitialDataFromFields(fields, sharedFieldsIds),
			options: {
				id: nanoid(6)
			}
		});

		this.currentData = fromStore(this.form.form);

		//

		watch(sharedData, (newSharedData) => {
			const dataOfNotOverriddenFields = pipe(
				newSharedData,
				// Only fields that are in the fields array
				R.filter((_, key) => this.fields.map((f) => f.CredimiID).includes(key)),
				// Not overridden
				R.filter((_, key) => !this.overriddenFieldsIds.includes(key))
			);

			this.form.form.update((oldData) => {
				return { ...oldData, ...dataOfNotOverriddenFields };
			});
		});

		//

		const formData = fromStore(this.form.form);

		watch(
			() => formData.current,
			() => {
				this.form.validateForm({ update: false }).then((result) => {
					console.log(result);
					this.isValid = result.valid;
				});
			}
		);
	}

	overrideField(id: string) {
		this.overriddenFieldsIds.push(id);
	}

	resetOverride(id: string) {
		this.overriddenFieldsIds = this.overriddenFieldsIds.filter((f) => f !== id);
		this.form.form.update((oldData) => {
			return { ...oldData, [id]: this.sharedData()[id] };
		});
	}
}

// export class TestConfigJsonForm {
// 	private readonly FORM_KEY = 'jsonConfig';
// 	readonly form: SuperForm<GenericRecord>;

// 	constructor(readonly jsonConfig: GenericRecord, readonly onValidUpdate: (jsonConfig: GenericRecord) => void) {
// 		this.form = createForm({
// 			adapter: zod(z.object({
// 				[this.FORM_KEY]: stringifiedObjectSchema.optional()
// 			})),
// 			initialData: {
// 				[this.FORM_KEY]: JSON.stringify(jsonConfig, null, 4)
// 			},
// 			options: {
// 				id: nanoid(6)
// 			}
// 		});

// 		const taintedState = fromStore(this.form.tainted);
// 		const isJsonConfigTainted = $derived(Boolean(taintedState.current?.jsonConfig));

// 		watch(isJsonConfigTainted, (isTainted) => {
// 			if (isTainted) {
// 				this.onValidUpdate(this.form.data.jsonConfig);
// 			}
// 		});
// 	}
// }
