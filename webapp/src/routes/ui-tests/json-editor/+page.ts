import { pb } from '@/pocketbase';
import { z } from 'zod';
import type { StandardWithTestSuites } from './logic';

//

const variantsSchema = z.object({
	variants: z.array(z.string())
});

export const load = async () => {
	const variantsResponse = await pb.send(
		'/api/conformance-checks/configs/get-configs-templates',
		{
			method: 'GET'
		}
	);
	const parseResult = variantsSchema.safeParse(variantsResponse);

	if (!parseResult.success) {
		throw new Error('Failed to parse variants response');
	}

	const standardsAndTestSuites: StandardWithTestSuites[] = [
		{
			id: 'openid4vp_wallet',
			label: 'OpenID4VP Wallet',
			description:
				'Lorem ipsum dolor sit amet consectetur. Tortor phasellus a feugiat mattis massa sollicitudin bibendum.',
			testSuites: [
				{
					id: 'openid-foundation-wallet',
					label: 'OpenID Foundation OpenID4VP Wallet',
					tests: parseResult.data.variants
				}
			]
		},
		{
			id: 'openid4vp_verifier',
			label: 'OpenID4VP Verifier',
			testSuites: [],
			description:
				'Lorem ipsum dolor sit amet consectetur. Tortor phasellus a feugiat mattis massa sollicitudin bibendum.'
		}
	];

	return {
		standardsAndTestSuites
	};
};
