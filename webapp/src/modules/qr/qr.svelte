<script lang="ts">
	import { generateQrCode } from './qr';
	import type { HTMLImgAttributes } from 'svelte/elements';

	type Props = HTMLImgAttributes & {
		src: string;
		cellSize?: number;
	};

	let { src, cellSize, ...rest }: Props = $props();

	const qr = generateQrCode(src, cellSize ?? 20);
</script>

<img src={qr} {...rest} />
