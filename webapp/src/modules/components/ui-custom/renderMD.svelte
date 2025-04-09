<!--
SPDX-FileCopyrightText: 2025 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { marked } from 'marked';

	import HTML from '@/components/ui-custom/renderHTML.svelte';

	type Props = {
		content: string;
		className?: string;
	};

	const { content, className }: Props = $props();
</script>

<div>
	<HTML {className} content={marked.parse(content, { async: false })} />
</div>
