<script lang="ts">
	import { page } from '$app/stores';
	import { A, Heading, P } from 'flowbite-svelte';
	const status = $page.status;
	const title = status === 404 ? 'Not Found' : 'Internal Error';
	const image = status === 404 ? '/404-computer.svg' : '/500.svg';
	const message = $page.error?.message;
</script>

<div
	class="mx-auto flex min-h-screen max-w-4xl flex-col gap-4 p-8 pt-20 lg:flex-row lg:p-0 lg:pt-0"
>
	<div class="flex w-full flex-col items-center justify-center">
		<img src={image} alt={title} />
	</div>
	<div class="flex w-full flex-col items-center justify-center gap-4">
		<Heading tag="h3" class="text-primary-600">{status} {title}</Heading>
		<Heading tag="h2">{message}</Heading>
		<div class="w-full pt-8">
			<P class="text-gray-400">Here are some Helpful link:</P>
			<ul class="flex gap-2">
				<li><A href="/">Home</A></li>
				<li><A href="/login">Login</A></li>
				<li><A href="/register">Register</A></li>
				<li><A href="/my">My</A></li>
			</ul>
		</div>
	</div>
</div>
