<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import Card from '@/components/ui-custom/card.svelte';

	interface Props {
		class?: string;
		contentClass?: string;
		children?: import('svelte').Snippet;
	}

	let { class: className = '', contentClass = '', children }: Props = $props();
</script>

<Card class="shadow-md ${className}" contentClass="space-y-6 {contentClass}">
	{@render children?.()}
</Card>
