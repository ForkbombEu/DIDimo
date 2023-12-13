<script>
	import { portal } from './utils/portal';

	export let target = 'body';
</script>

<div use:portal={target} class="fixed m-0 p-0 z-50">
	<slot />
</div>
