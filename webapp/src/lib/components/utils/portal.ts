import { tick } from 'svelte';

// Credits: https://github.com/romkor/svelte-portal/blob/master/src/Portal.svelte
export function portal(
	el: HTMLElement,
	target: HTMLElement | string = 'body > div'
): {
	update: (newTarget: HTMLElement | string) => Promise<void>;
	destroy: () => void;
} {
	let targetEl: HTMLElement | null = null;

	async function update(newTarget: HTMLElement | string): Promise<void> {
		target = newTarget;
		if (typeof target === 'string') {
			targetEl = targetEl ?? document.querySelector(target);

			if (!targetEl) {
				await tick();
				targetEl = document.querySelector(target);
			}

			if (!targetEl) throw new Error(`No element found matching CSS selector: "${target}"`);
		} else if (target instanceof HTMLElement) {
			targetEl = target;
		} else {
			throw new TypeError(
				`Unknown portal target type: ${
					target === null ? 'null' : typeof target
				}. Allowed types: string (CSS selector) or HTMLElement.`
			);
		}

		targetEl.appendChild(el);
		el.hidden = false;
	}

	function destroy(): void {
		if (el.parentElement) el.parentElement.removeChild(el);
	}

	update(target);

	return { update, destroy };
}
