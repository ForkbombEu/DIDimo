import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';
import type { UsersResponse } from './types';

export const pb = new PocketBase(PUBLIC_POCKETBASE_URL);
export const currentUser = writable(pb.authStore.model as AuthStoreModel);

export type AuthStoreModel = UsersResponse | null;
