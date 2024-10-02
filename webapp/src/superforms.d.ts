// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

export * from 'sveltekit-superforms/client';

declare module 'sveltekit-superforms/client' {
	function formFieldProxy<
		T extends ZodValidation<z.AnyZodObject>,
		Path extends FormPathLeaves<z.TypeOf<UnwrapEffects<T>>>
	>(
		form: SuperForm<T, import('./lib/errorHandling').ClientResponseErrorData>,
		path: Path
	): {
		path: Path;
		value: Writable<FormPathType<z.infer<UnwrapEffects<T>>, Path>>;
		errors: Writable<string[] | undefined>;
		constraints: Writable<InputConstraint | undefined>;
	};
}
