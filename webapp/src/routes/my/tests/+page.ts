import { pb } from '@/pocketbase';
import { z } from 'zod';
import type { StandardWithTestSuites } from './logic';

//

const variantsSchema = z.object({
	variants: z.array(z.string())
});

export const load = async ({ fetch }) => {
	const openid4vpVariantsResponse = await pb.send(
		'/api/conformance-checks/configs/get-configs-templates',
		{
			method: 'GET',
			fetch
		}
	);

	const openid4vciVariantsResponse = await pb.send(
		'/api/conformance-checks/configs/get-configs-templates?test_id=OpenID4VCI_Wallet/eudiw',
		{
			method: 'GET',
			fetch
		}
	);

	const parseOpenid4vpResult = variantsSchema.safeParse(openid4vpVariantsResponse);
	const parseOpenid4vciResult = variantsSchema.safeParse(openid4vciVariantsResponse);

	if (!parseOpenid4vpResult.success || !parseOpenid4vciResult.success) {
		throw new Error('Failed to parse variants response');
	}

	const standardsAndTestSuites: StandardWithTestSuites[] = [
		{
			id: 'OpenID4VP_Wallet',
			label: 'OpenID4VP Wallet',
			description:
				'Lorem ipsum dolor sit amet consectetur. Tortor phasellus a feugiat mattis massa sollicitudin bibendum.',
			testSuites: [
				{
					id: 'OpenID_Foundation',
					label: 'OpenID4VP Wallet',
					tests: parseOpenid4vpResult.data.variants
				}
			]
		},
		{
			id: 'OpenID4VCI_Wallet',
			label: 'OpenID4VCI Wallet',
			description:
				'Lorem ipsum dolor sit amet consectetur. Tortor phasellus a feugiat mattis massa sollicitudin bibendum.',
			testSuites: [
				{
					id: 'eudiw',
					label: 'OpenID4VCI Wallet',
					tests: parseOpenid4vciResult.data.variants
				}
			]
		}
	];

	return {
		standardsAndTestSuites
	};
};
