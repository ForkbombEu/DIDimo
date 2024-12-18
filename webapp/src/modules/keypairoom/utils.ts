import { pb } from '@/pocketbase';
import {
	Collections,
	type UsersPublicKeysRecord,
	type UsersPublicKeysResponse
} from '@/pocketbase/types';
import _ from 'lodash';
import type { Keypair } from './keypair';
import { createSessionStorageHandlers } from '@/utils/sessionStorage';

export type PublicKeys = Omit<UsersPublicKeysRecord, 'owner'>;

export function getPublicKeysFromKeypair(keypair: Keypair): PublicKeys {
	const publicKeys = _.cloneDeep(keypair);
	// @ts-expect-error Cannot use delete on required field
	delete publicKeys.seed;
	// @ts-expect-error Cannot use delete on required field
	delete publicKeys.keyring;
	return publicKeys;
}

export async function getUserPublicKeys(userId: string | undefined = undefined, fetchFn = fetch) {
	const id = userId ?? pb.authStore.model?.id ?? '';
	try {
		return await pb
			.collection(Collections.UsersPublicKeys)
			.getFirstListItem<UsersPublicKeysResponse>(`owner.id = '${id}'`, { fetch: fetchFn });
	} catch (e) {
		return undefined;
	}
}

export async function saveUserPublicKeys(publicKeys: PublicKeys) {
	const data: UsersPublicKeysRecord = {
		...publicKeys,
		owner: pb.authStore.model?.id
	};
	await pb.collection('users_public_keys').create(data);
}

//

export const RegenerateKeyringSession = createSessionStorageHandlers('keypairoom_regenerate');
