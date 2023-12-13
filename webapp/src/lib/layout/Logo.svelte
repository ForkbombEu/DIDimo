<script>
	import { appTitle } from '$lib/strings';
	import { NavBrand } from 'flowbite-svelte';

	/*
	 * The url to navigate to when the logo is clicked.
	 */
	export let href = '/my';
</script>

<NavBrand {href}>
	<img src={'/logo.svg'} class="h-9" alt={appTitle} />
</NavBrand>

<!--  
    @component 
    Simple logo component. Change logo.svg to your logo.

    ### Props
    - `href` - The url to navigate to when the logo is clicked.
-->
