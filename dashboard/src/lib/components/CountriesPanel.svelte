<script>
	import { List, Globe2 } from 'lucide-svelte';
	import TopItemsList from './TopItemsList.svelte';
	import GeographicalMap from './GeographicalMap.svelte';
	import { Button } from '$lib/components/ui/button';

	let { countries = [], onclick = null, loading = false } = $props();

	let activeTab = $state('map'); // Default to 'map' view

	function handleCountryClick(item) {
		if (onclick) {
			onclick(item);
		}
	}
</script>

<div class="space-y-4">
	<!-- Tabs -->
	<div class="border-border flex gap-2 border-b">
		<Button
			variant={activeTab === 'list' ? 'default' : 'ghost'}
			size="sm"
			class="rounded-b-none"
			onclick={() => (activeTab = 'list')}
		>
			<List class="mr-2 h-4 w-4" />
			List View
		</Button>
		<Button
			variant={activeTab === 'map' ? 'default' : 'ghost'}
			size="sm"
			class="rounded-b-none"
			onclick={() => (activeTab = 'map')}
		>
			<Globe2 class="mr-2 h-4 w-4" />
			Map View
		</Button>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
			<div class="flex flex-col items-center gap-2">
				<div
					class="border-primary h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
				></div>
				<p class="text-sm">Loading countries...</p>
			</div>
		</div>
	{:else if activeTab === 'list'}
		<TopItemsList
			items={countries}
			labelKey="name"
			valueKey="count"
			maxItems={5}
			type="country"
			showMoreTitle="All Countries ({countries.length} total)"
			{onclick}
		/>
	{:else}
		<GeographicalMap
			data={countries.map((c) => ({ country: c.name, count: c.count }))}
			onclick={handleCountryClick}
		/>
	{/if}
</div>
