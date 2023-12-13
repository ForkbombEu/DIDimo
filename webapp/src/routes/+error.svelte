<script lang="ts">
	import { goto } from '$app/navigation';
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
	class="flex flex-col gap-4 lg:flex-row max-w-4xl mx-auto min-h-screen pt-20 lg:pt-0 p-8 lg:p-0"
>
	<div class="flex flex-col items-center justify-center w-full">
		<img src={data.image} alt={data.title} />
	</div>
	<div class="flex flex-col items-center justify-center w-full gap-4">
		<Heading tag="h3" class="text-primary-600">{data.status} {data.title}</Heading>
		<Heading tag="h2">{data.description}</Heading>
		{#if status !== 503}
			<div class="w-full pt-8 space-y-1">
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
