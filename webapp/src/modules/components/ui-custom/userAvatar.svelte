<script lang="ts">
	import Avatar from './avatar.svelte';
	import type { UsersResponse } from '@/pocketbase/types';
	import { pb } from '@/pocketbase';
	import { m } from '@/i18n';
	import type { ComponentProps } from 'svelte';

	type Props = ComponentProps<typeof Avatar> & { user?: UsersResponse };

	let { user = pb.authStore.model as UsersResponse, ...rest }: Props = $props();
	if (!user) throw new Error('missing user');

	const src = $derived(pb.files.getURL(user, user.avatar));
	const fallback = $derived(user.name.slice(0, 2));
</script>

<Avatar {...rest} {src} {fallback} alt="{m.Avatar()} {user.name}" />
