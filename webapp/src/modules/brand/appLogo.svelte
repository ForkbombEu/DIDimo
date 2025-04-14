<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { appName } from '@/brand';
	import { localizeHref } from '@/i18n';

	interface Props {
		/*
		 * The url to navigate to when the logo is clicked.
		 */
		href?: string;
	}

	let { href = '/' }: Props = $props();
</script>

<a href={localizeHref(href)}>
	<img src={'/logos/credimi_logo-transp_emblem.png'} class="h-9 w-auto" alt={appName} />
</a>

<!--  
    @component 
    Simple logo component. Change logo.svg to your logo.

    ### Props
    - `href` - The url to navigate to when the logo is clicked.
-->
