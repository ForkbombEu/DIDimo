<script lang="ts">
	import { appName } from '@/brand';

	interface Props {
		/*
		 * The url to navigate to when the logo is clicked.
		 */
		href?: string;
	}

	let { href = '/' }: Props = $props();
</script>

<a {href}>
	<img src={'/logo.svg'} class="h-9" alt={appName} />
</a>

<!--  
    @component 
    Simple logo component. Change logo.svg to your logo.

    ### Props
    - `href` - The url to navigate to when the logo is clicked.
-->
