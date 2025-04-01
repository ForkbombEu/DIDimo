import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';
import type { TypedPocketBase, UsersResponse } from '@/pocketbase/types';

//

export const pb = new PocketBase() as TypedPocketBase;

export const currentUser = writable(pb.authStore.model as AuthStoreModel);
export type AuthStoreModel = UsersResponse | null;
