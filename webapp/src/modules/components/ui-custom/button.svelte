<script lang="ts">
	import { cn } from '@/components/ui/utils';
	import type { ButtonProps } from '@/components/ui/button';
	import { Button } from '@/components/ui/button';
	import { localizeHref } from '@/i18n/paraglide/runtime';

	let { class: className, href, ...rest }: ButtonProps = $props();
</script>

<Button href={href ? localizeHref(href) : undefined} class={cn(className)} {...rest} />
