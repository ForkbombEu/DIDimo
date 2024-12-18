<script lang="ts">
	import Avatar, { type AvatarProps } from '@/components/ui-custom/avatar.svelte';
	import { cn } from '@/components/ui/utils';
	import { pb } from '@/pocketbase';
	import type { OrganizationsRecord } from '@/pocketbase/types';

	type Props = AvatarProps & {
		organization: OrganizationsRecord;
	};

	let { organization, ...rest }: Props = $props();

	let src = $derived(pb.files.getUrl(organization, organization.avatar ?? ''));
	let fallback = $derived(organization.name.slice(0, 2));
</script>

<Avatar {...rest} {src} {fallback} class={cn(rest.class, 'rounded-sm')} />
