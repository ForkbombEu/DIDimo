<script lang="ts">
	import { appName } from '@/brand';
	import A from '@/components/ui-custom/a.svelte';

	interface Props {
		/*
		 * The url to navigate to when the logo is clicked.
		 */
		href?: string;
	}

	let { href = '/' }: Props = $props();
</script>

<A {href}>
	<img src={'/logo.svg'} class="h-9" alt={appName} />
</A>

<!--  
    @component 
    Simple logo component. Change logo.svg to your logo.

    ### Props
    - `href` - The url to navigate to when the logo is clicked.
-->
