import { watch } from 'runed';
import type { z } from 'zod';

export interface FieldData {
	id: string;
	schema: z.ZodTypeAny;
	label: string;
	description: string;
}

export class Field implements FieldData {
	constructor(
		public id: string,
		public defaultValue: () => unknown,
		public schema: z.ZodTypeAny,
		//
		public label: string,
		public description: string
	) {
		this.value = defaultValue();

		watch(defaultValue, (v) => {
			if (!this.isOverridden) this.value = v;
		});
	}

	value = $state<unknown>();
	status = $derived.by(() => this.schema.safeParse(this.value));
	isValid = $derived.by(() => this.status.success);
	error = $derived.by(() => this.status.error?.message);
	hasError = $derived.by(() => this.status.error);

	isOverridden = $state(false);

	resetOverride() {
		this.isOverridden = false;
		this.value = this.defaultValue();
	}

	getSchemaType() {
		if (this.schema._def.typeName == 'ZodString') return 'text';
		if (this.schema._def.typeName == 'ZodNumber') return 'number';
		if (this.schema._def.typeName == 'ZodBoolean') return 'boolean';
		return 'unknown';
	}
}

export class Form {
	fields: Field[] = [];
	onSubmit: (data: Record<string, unknown>) => void | Promise<void>;

	constructor(config: {
		fieldsData: FieldData[];
		defaultValues: () => Record<string, unknown>;
		onSubmit: (data: Record<string, unknown>) => void | Promise<void>;
	}) {
		this.onSubmit = config.onSubmit;

		this.fields = config.fieldsData.map(
			(f) =>
				new Field(
					f.id,
					() => config.defaultValues()[f.id],
					f.schema,
					f.label,
					f.description
				)
		);
	}

	isValid = $derived.by(() => {
		const statuses = this.fields.map((f) => f.status);
		return statuses.every((s) => s.success);
	});

	validData = $derived.by(() => {
		return Object.fromEntries(
			this.fields.map((f) => [f.id, f.status.success ? f.value : undefined])
		);
	});

	errors = $derived.by(() => {
		return Object.fromEntries(
			this.fields.map((f) => [f.id, f.status.error ? f.status.error.message : undefined])
		);
	});

	submit() {
		if (!this.isValid) return;
		this.onSubmit(this.validData);
	}
}
