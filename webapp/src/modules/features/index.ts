import { Features } from './features-list.generated';
import { pb } from '@/pocketbase';
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

//

export const featureFlags = writable<FeatureFlags>();

export type FeatureKeys = keyof typeof Features;
export type FeatureFlags = Record<FeatureKeys, boolean>;

export async function loadFeatureFlags(fetchFn = fetch): Promise<FeatureFlags> {
	const flags: Partial<FeatureFlags> = {};

	const list = await pb.collection('flags').getFullList({ requestKey: null, fetch: fetchFn });

	for (const [key, name] of Object.entries(Features)) {
		const feature = list.find((f) => f.name === name);
		flags[key as FeatureKeys] = feature?.active ?? false;
	}

	if (browser) featureFlags.set(flags as FeatureFlags);
	return flags as FeatureFlags;
}

export * from './features-list.generated';
