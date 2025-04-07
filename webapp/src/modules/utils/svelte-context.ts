// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { getContext as svelteGetContext, setContext as svelteSetContext } from 'svelte';

//

export function setupDerivedContext<Context>(key: string) {
	const contextKey = Symbol(key);

	function setDerivedContext<CustomContext = Context>(derived: () => CustomContext) {
		return svelteSetContext(contextKey, derived);
	}

	/**
	 * @returns a value that can be `$derived(...)` in order to get a reactive context
	 */
	function getDerivedContext<CustomContext = Context>() {
		const baseContext =
			svelteGetContext<ReturnType<typeof setDerivedContext<CustomContext>>>(contextKey);
		return baseContext();
	}

	return {
		setDerivedContext,
		getDerivedContext,
		contextKey
	};
}
