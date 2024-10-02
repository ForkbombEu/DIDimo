<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { Navbar } from 'flowbite-svelte';
</script>

<div class="shrink-0 border-b">
	<Navbar fluid>
		<div class="shrink-0 items-center justify-start">
			<slot name="left" />
		</div>
		<div class="justify-center">
			<slot name="center" />
		</div>
		<div class="shrink-0 flex items-center justify-end">
			<slot name="right" />
		</div>
	</Navbar>
</div>

<!-- 
	@component 

	### Slots
	- left - the left side of the navbar
	- center - the center of the navbar
	- right - the right side of the navbar
	
-->
