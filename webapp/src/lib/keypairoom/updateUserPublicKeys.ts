import { pb } from '$lib/pocketbase';
import { Collections, type UsersRecord } from '$lib/pocketbase/types';
import { log } from '$lib/utils/devLog';
import type { Keypair } from './keypair';

export type UserPublicKeys = Omit<UsersRecord, 'name' | 'avatar'>;

export function getPublicKeysFromKeypair(keypair: Partial<Keypair>): UserPublicKeys {
	delete keypair.seed;
	delete keypair.keyring;
	return keypair;
}

export async function updateUserPublicKeys(userId: string, publicKeys: UserPublicKeys) {
	try {
		await pb.collection(Collections.Users).update(userId, publicKeys);
	} catch (e) {
		log(e);
	}
}
