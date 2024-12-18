<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { Moon } from 'svelte-loading-spinners';

	type Props = ComponentProps<Moon>;

	const { size = 40, color = 'black', ...rest }: Props = $props();
</script>

<Moon {size} {color} {...rest} />
