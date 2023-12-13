import type { FormPathLeaves, ZodValidation } from 'sveltekit-superforms';
import type { SuperForm } from 'sveltekit-superforms/client';
import type { AnyZodObject, z } from 'zod';

export interface FormComponentProps<T extends AnyZodObject, O extends Record<string, unknown>> {
	superform: SuperForm<ZodValidation<T>, any>;
	field: FormPathLeaves<z.infer<T>>;
	options: O;
}
