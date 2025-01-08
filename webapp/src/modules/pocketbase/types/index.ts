export * from './index.generated';
export * from './extra.generated';

//

import type { SimplifyDeep } from 'type-fest';
export type Data<R extends Record<string, unknown>> = SimplifyDeep<
	Omit<R, 'id' | 'created' | 'updated'>
>;
