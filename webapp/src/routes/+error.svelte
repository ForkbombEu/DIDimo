<!--
SPDX-FileCopyrightText: 2024 The Forkbomb Company

SPDX-License-Identifier: AGPL-3.0-or-later
-->

<script lang="ts">
	import { page } from '$app/stores';
	import { A, Heading, P } from 'flowbite-svelte';
	type ErrorData = {
		status: number;
		image: string;
		title: string;
		description: string;
	};
	const status = $page.status;
	const data: ErrorData = (() => {
		switch (status) {
			case 404:
				return {
					status: 404,
					image: '/404-computer.svg',
					title: 'Not Found',
					description: 'Whoops! That page doesn’t exist.'
				};
			case 503:
				return {
					status: 503,
					image: '/maintenance.svg',
					title: 'Service Unavailable',
					description: 'We’re temporarily offline for maintenance. Please try again later.'
				};
			default:
				return {
					status: 500,
					image: '/500.svg',
					title: 'Internal Error',
					description: 'Something went wrong. Please try again later.'
				};
		}
	})();
</script>

<div
	class="mx-auto flex min-h-screen max-w-4xl flex-col gap-4 p-8 pt-20 lg:flex-row lg:p-0 lg:pt-0"
>
	<div class="flex w-full flex-col items-center justify-center">
		<img src={data.image} alt={data.title} />
	</div>
	<div class="flex w-full flex-col items-center justify-center gap-4">
		<Heading tag="h3" class="text-primary-600">{data.status} {data.title}</Heading>
		<Heading tag="h2">{data.description}</Heading>
		{#if status !== 503}
			<div class="w-full space-y-1 pt-8">
				<P class="text-gray-400">Here are some Helpful link:</P>
				<ul class="flex gap-6">
					<li><A href="/">Home</A></li>
					<li><A href="/login">Login</A></li>
					<li><A href="/register">Register</A></li>
				</ul>
			</div>
		{/if}
	</div>
</div>
