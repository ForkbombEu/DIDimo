<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script>
	import UserAvatar from '$lib/components/userAvatar.svelte';
	import { Dropdown } from 'flowbite-svelte';
</script>

<button id="avatar-menu">
	<UserAvatar />
</button>
<Dropdown
	placement="bottom"
	triggeredBy="#avatar-menu"
	footerClass="w-full min-h-screen md:w-fit md:min-h-fit"
>
	<slot />
</Dropdown>
