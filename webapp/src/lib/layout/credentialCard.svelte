<script lang="ts">
	import T from '@/components/ui-custom/t.svelte';
	import { Card } from '@/components/ui/card';
	import { m } from '@/i18n';

	type Credential = {
		name: string;
		issuer: string;
		duration: string;
		category: string;
		format: string;
		logo?: string;
	};

	type Props = {
		credential?: Credential;
		class?: string;
	};

	const { credential = defaultCredential(), class: className = '' }: Props = $props();

	function defaultCredential(): Credential {
		return {
			name: "Driver's License",
			issuer: 'Didimo',
			duration: '1 year',
			category: 'Identity',
			format: 'Digital',

			logo:
				Math.random() > 0.5
					? 'https://upload.wikimedia.org/wikipedia/commons/thumb/2/24/LEGO_logo.svg/2048px-LEGO_logo.svg.png'
					: undefined
		};
	}
</script>

<Card class="border-primary flex flex-col justify-between space-y-6 rounded-xl p-6 {className}">
	<div class="flex items-center space-x-4">
		{#if credential.logo}
			<img src={credential.logo} alt={credential.name} class="size-14 rounded-lg border" />
		{:else}
			<div class="flex size-14 items-center justify-center rounded-lg border bg-gray-100">
				<p class="font-semibold uppercase">{credential.name.slice(0, 2)}</p>
			</div>
		{/if}
		<T tag="h4">{credential.name}</T>
	</div>

	<div class="grid grid-cols-2">
		<div>
			<small>{m.Issuer()}:</small>
			<small class="text-muted-foreground font-semibold">{credential.issuer}</small>
		</div>
		<div>
			<small>{m.Duration()}:</small>
			<small class="text-muted-foreground font-semibold">{credential.duration}</small>
		</div>
		<div>
			<small>{m.Category()}:</small>
			<small class="text-muted-foreground font-semibold">{credential.category}</small>
		</div>
		<div>
			<small>{m.Format()}:</small>
			<small class="text-muted-foreground font-semibold">{credential.format}</small>
		</div>
	</div>
</Card>
