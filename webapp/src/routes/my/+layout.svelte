<script lang="ts">
	import * as Sidebar from '@/components/ui/sidebar/index.js';
	import MySidebar from './_partials/mySidebar.svelte';
	import type { Snippet } from 'svelte';
	import SidebarContext from '@/components/layout/sidebar/sidebarContext.svelte';
	import Topbar from '@/components/layout/topbar.svelte';
	import { AppLogo } from '@/brand';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();
</script>

<Sidebar.Provider>
	<MySidebar />

	<main class="grow">
		<SidebarContext>
			{#snippet content({ isMobile })}
				{#if isMobile}
					<Topbar>
						{#snippet left()}
							<div class="flex items-center gap-2">
								<Sidebar.Trigger variant="outline" class="size-9" />
								<AppLogo />
							</div>
						{/snippet}
						{#snippet right()}{/snippet}
					</Topbar>
				{/if}
			{/snippet}
		</SidebarContext>

		{@render children?.()}
	</main>
</Sidebar.Provider>
