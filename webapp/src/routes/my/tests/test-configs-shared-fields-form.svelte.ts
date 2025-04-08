// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import type { GenericRecord } from '@/utils/types';
import { createInitialDataFromFields, createTestConfigFormSchema, type ConfigField } from './logic';
import type { SuperForm } from 'sveltekit-superforms';
import { createForm } from '@/forms';
import { zod } from 'sveltekit-superforms/adapters';
import { nanoid } from 'nanoid';
import { fromStore } from 'svelte/store';
import { watch } from 'runed';
import { Record } from 'effect';

//

export class TestConfigsSharedFieldsForm {
	readonly form: SuperForm<GenericRecord>;

	constructor(
		readonly fields: ConfigField[],
		readonly onUpdate: (data: GenericRecord) => void
	) {
		this.form = createForm({
			adapter: zod(createTestConfigFormSchema(fields)),
			initialData: createInitialDataFromFields(fields),
			options: {
				id: nanoid(6)
			}
		});

		const { form: formData, validateForm } = this.form;
		const formState = fromStore(formData);

		watch(
			() => formState.current,
			(updatedData) => {
				validateForm({ update: false }).then((result) => {
					const invalidFieldsIds = Object.keys(result.errors);
					const validData = Record.filter(
						updatedData,
						(_, id) => !invalidFieldsIds.includes(id)
					);
					this.onUpdate(validData);
				});
			}
		);
	}
}
