// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { getUserPublicKeys } from '@/keypairoom/utils';

export async function getUserDidUrl(userId: string | undefined = undefined) {
	const publicKeys = await getUserPublicKeys(userId);
	if (publicKeys) return createDidUrl(publicKeys.eddsa_public_key);
}

export function createDidUrl(eddsaPublicKey: string) {
	return `https://explorer.did.dyne.org/details/did:dyne:sandbox.signroom:${eddsaPublicKey}`;
}
