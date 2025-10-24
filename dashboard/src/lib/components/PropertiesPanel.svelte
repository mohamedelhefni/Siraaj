<script>
	import { Search, Filter, X, ChevronDown } from 'lucide-svelte';
	import { Badge } from '$lib/components/ui/badge';
	import Modal from './Modal.svelte';

	let { properties = [], onPropertyClick = null } = $props();

	let searchTerm = $state('');
	let selectedKey = $state(null);
	let modalOpen = $state(false);
	let showMoreModalOpen = $state(false);
	let selectedProperty = $state(null);
	let maxDisplayed = $state(5); // Show only 5 properties initially

	// Filter properties based on search and selected key
	const filteredProperties = $derived(
		properties.filter((prop) => {
			const matchesSearch =
				searchTerm === '' ||
				prop.key.toLowerCase().includes(searchTerm.toLowerCase()) ||
				prop.value.toLowerCase().includes(searchTerm.toLowerCase());
			const matchesKey = !selectedKey || prop.key === selectedKey;
			return matchesSearch && matchesKey;
		})
	);

	// Properties to display (limited)
	const displayedProperties = $derived(filteredProperties.slice(0, maxDisplayed));
	const hasMore = $derived(filteredProperties.length > maxDisplayed);

	// Get unique keys for filtering
	const uniqueKeys = $derived([...new Set(properties.map((p) => p.key))].sort());

	// Get total count
	const totalCount = $derived(filteredProperties.reduce((sum, prop) => sum + (prop.count || 0), 0));

	function handlePropertyClick(prop) {
		if (onPropertyClick) {
			onPropertyClick(prop);
		}
	}

	function viewDetails(prop) {
		selectedProperty = prop;
		modalOpen = true;
	}

	function clearKeyFilter() {
		selectedKey = null;
	}

	function openShowMore() {
		showMoreModalOpen = true;
	}
</script>

<div class="space-y-4">
	<!-- Search and Filter Bar -->
	<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
		<div class="relative flex-1">
			<Search class="text-muted-foreground absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2" />
			<input
				type="text"
				bind:value={searchTerm}
				placeholder="Search properties..."
				class="border-input bg-background focus:ring-ring h-10 w-full rounded-md border pl-10 pr-3 text-sm focus:outline-none focus:ring-2"
			/>
		</div>

		{#if selectedKey}
			<Badge variant="secondary" class="flex w-fit items-center gap-2">
				<Filter class="h-3 w-3" />
				{selectedKey}
				<button
					type="button"
					onclick={clearKeyFilter}
					class="hover:bg-accent ml-1 rounded-full p-0.5"
				>
					<X class="h-3 w-3" />
				</button>
			</Badge>
		{/if}
	</div>

	<!-- Key Filter Pills -->
	{#if uniqueKeys.length > 0}
		<div class="flex flex-wrap gap-2">
			{#each uniqueKeys.slice(0, 10) as key}
				<button
					type="button"
					onclick={() => (selectedKey = selectedKey === key ? null : key)}
					class="hover:bg-accent rounded-full px-3 py-1 text-xs font-medium transition-colors {selectedKey ===
					key
						? 'bg-primary text-primary-foreground'
						: 'bg-secondary text-secondary-foreground'}"
				>
					{key}
				</button>
			{/each}
		</div>
	{/if}

	<!-- Properties List -->
	<div class="space-y-2">
		{#if filteredProperties.length === 0}
			<p class="text-muted-foreground py-8 text-center text-sm">
				{#if searchTerm || selectedKey}
					No properties match your filters
				{:else}
					No properties available
				{/if}
			</p>
		{:else}
			<div class="text-muted-foreground mb-2 text-xs">
				Showing {displayedProperties.length} of {filteredProperties.length} propert{filteredProperties.length ===
				1
					? 'y'
					: 'ies'}
				· Total events: {totalCount.toLocaleString()}
			</div>

			{#each displayedProperties as prop}
				{@const percentage = totalCount > 0 ? ((prop.count / totalCount) * 100).toFixed(1) : 0}
				<div
					class="hover:bg-accent group relative flex items-center justify-between rounded-lg border p-3 transition-colors"
				>
					<div class="min-w-0 flex-1">
						<div class="mb-1 flex items-center gap-2">
							<code
								class="bg-muted text-muted-foreground truncate rounded px-2 py-0.5 text-xs font-medium"
							>
								{prop.key}
							</code>
							<span class="text-primary truncate text-sm font-medium" title={prop.value}>
								{prop.value}
							</span>
						</div>
						<div class="flex items-center gap-3 text-xs">
							<span class="text-muted-foreground">
								{prop.count?.toLocaleString()} events ({percentage}%)
							</span>
							{#if prop.event_types > 1}
								<Badge variant="outline" class="text-xs">
									{prop.event_types} event types
								</Badge>
							{/if}
						</div>
					</div>

					<div class="ml-4 flex gap-2 opacity-0 transition-opacity group-hover:opacity-100">
						<button
							type="button"
							onclick={() => viewDetails(prop)}
							class="hover:bg-background rounded-md px-3 py-1 text-xs font-medium transition-colors"
						>
							Details
						</button>
						{#if onPropertyClick}
							<button
								type="button"
								onclick={() => handlePropertyClick(prop)}
								class="bg-primary text-primary-foreground hover:bg-primary/90 rounded-md px-3 py-1 text-xs font-medium transition-colors"
							>
								Filter
							</button>
						{/if}
					</div>
				</div>
			{/each}

			<!-- Show More Button -->
			{#if hasMore}
				<button
					type="button"
					onclick={openShowMore}
					class="hover:bg-accent text-muted-foreground flex w-full items-center justify-center gap-2 rounded-lg border border-dashed p-3 text-sm font-medium transition-colors"
				>
					<ChevronDown class="h-4 w-4" />
					Show All ({filteredProperties.length} total)
				</button>
			{/if}
		{/if}
	</div>
</div>

<!-- Property Details Modal -->
<Modal bind:open={modalOpen} title="Property Details">
	{#snippet children()}
		{#if selectedProperty}
			<div class="space-y-4">
				<div>
					<h3 class="text-muted-foreground mb-1 text-xs font-medium uppercase">Key</h3>
					<code class="bg-muted block rounded-md p-3 text-sm">{selectedProperty.key}</code>
				</div>

				<div>
					<h3 class="text-muted-foreground mb-1 text-xs font-medium uppercase">Value</h3>
					<div class="bg-muted rounded-md p-3">
						<span class="break-all text-sm">{selectedProperty.value}</span>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<h3 class="text-muted-foreground mb-1 text-xs font-medium uppercase">Event Count</h3>
						<p class="text-2xl font-bold">{selectedProperty.count?.toLocaleString()}</p>
					</div>

					<div>
						<h3 class="text-muted-foreground mb-1 text-xs font-medium uppercase">Event Types</h3>
						<p class="text-2xl font-bold">{selectedProperty.event_types}</p>
					</div>
				</div>

				{#if onPropertyClick}
					<button
						type="button"
						onclick={() => {
							handlePropertyClick(selectedProperty);
							modalOpen = false;
						}}
						class="bg-primary text-primary-foreground hover:bg-primary/90 w-full rounded-md px-4 py-2 text-sm font-medium transition-colors"
					>
						Apply as Filter
					</button>
				{/if}
			</div>
		{/if}
	{/snippet}
</Modal>

<!-- Show More Modal -->
<Modal bind:open={showMoreModalOpen} title="All Properties ({filteredProperties.length})">
	{#snippet children()}
		<div class="max-h-[60vh] space-y-2 overflow-y-auto">
			{#each filteredProperties as prop}
				{@const percentage = totalCount > 0 ? ((prop.count / totalCount) * 100).toFixed(1) : 0}
				<div
					class="hover:bg-accent group relative flex items-center justify-between rounded-lg border p-3 transition-colors"
				>
					<div class="min-w-0 flex-1">
						<div class="mb-1 flex items-center gap-2">
							<code
								class="bg-muted text-muted-foreground truncate rounded px-2 py-0.5 text-xs font-medium"
							>
								{prop.key}
							</code>
							<span class="text-primary truncate text-sm font-medium" title={prop.value}>
								{prop.value}
							</span>
						</div>
						<div class="flex items-center gap-3 text-xs">
							<span class="text-muted-foreground">
								{prop.count?.toLocaleString()} events ({percentage}%)
							</span>
							{#if prop.event_types > 1}
								<Badge variant="outline" class="text-xs">
									{prop.event_types} event types
								</Badge>
							{/if}
						</div>
					</div>

					<div class="ml-4 flex gap-2 opacity-0 transition-opacity group-hover:opacity-100">
						<button
							type="button"
							onclick={() => {
								selectedProperty = prop;
								showMoreModalOpen = false;
								modalOpen = true;
							}}
							class="hover:bg-background rounded-md px-3 py-1 text-xs font-medium transition-colors"
						>
							Details
						</button>
						{#if onPropertyClick}
							<button
								type="button"
								onclick={() => {
									handlePropertyClick(prop);
									showMoreModalOpen = false;
								}}
								class="bg-primary text-primary-foreground hover:bg-primary/90 rounded-md px-3 py-1 text-xs font-medium transition-colors"
							>
								Filter
							</button>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/snippet}
</Modal>
