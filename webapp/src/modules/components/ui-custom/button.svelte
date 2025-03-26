<script lang="ts">
	import type { ButtonProps } from '@/components/ui/button';
	import Button from '@/components/ui-custom/button.svelte';
	import { localizeHref } from '@/i18n/paraglide/runtime';

	let { class: className, href, ...rest }: ButtonProps = $props();
</script>

<Button href={href ? localizeHref(href) : undefined} class={className} {...rest} />
