import { getExceptionMessage } from '@/utils/errors';
import type { GenericRecord } from '@/utils/types';
import type { FormOptions as SuperformOptions } from 'sveltekit-superforms';
import { type ValidationAdapter } from 'sveltekit-superforms/adapters';
import { defaults, setError, superForm } from 'sveltekit-superforms/client';

//

export type SubmitFunction<Data extends GenericRecord = GenericRecord> = NonNullable<
	SuperformOptions<Data>['onUpdate']
>;

export type FormOptions<Data extends GenericRecord = GenericRecord> = Omit<
	SuperformOptions<Data>,
	'onUpdate'
>;

export type CreateFormProps<Data extends GenericRecord> = {
	adapter: ValidationAdapter<Data>;
	options?: FormOptions<Data>;
	onSubmit?: SubmitFunction<Data>;
	initialData?: Partial<Data>;
};

export function createForm<Data extends GenericRecord>(props: CreateFormProps<Data>) {
	const { adapter, initialData = {} as Partial<Data>, options = {}, onSubmit = () => {} } = props;
	const form = defaults(initialData, adapter);
	return superForm(form, {
		SPA: true,
		applyAction: false,
		scrollToError: 'smooth',
		validators: adapter,
		dataType: 'json',
		taintedMessage: null,
		onUpdate: async (event) => {
			try {
				if (event.form.valid) await onSubmit(event);
			} catch (e) {
				setError(event.form, getExceptionMessage(e));
			}
		},
		...options
	});
}

export const FORM_ERROR_PATH = '_errors';
