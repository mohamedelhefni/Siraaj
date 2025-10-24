<script>
	let { items = [], labelKey = 'name', valueKey = 'count', maxItems = 10 } = $props();

	const displayItems = $derived(items.slice(0, maxItems));
	const total = $derived(items.reduce((sum, item) => sum + (item[valueKey] || 0), 0));
</script>

<div class="space-y-2">
	{#if displayItems.length === 0}
		<p class="text-sm text-muted-foreground text-center py-4">No data available</p>
	{:else}
		{#each displayItems as item, i}
			{@const percentage = total > 0 ? ((item[valueKey] / total) * 100).toFixed(1) : 0}
			<div class="flex items-center justify-between space-x-2">
				<div class="flex-1 min-w-0">
					<div class="flex items-center justify-between mb-1">
						<span class="text-sm font-medium truncate" title={item[labelKey]}>
							{i + 1}. {item[labelKey]}
						</span>
						<span class="text-sm text-muted-foreground ml-2">
							{item[valueKey]?.toLocaleString()} ({percentage}%)
						</span>
					</div>
					<div class="w-full bg-secondary rounded-full h-2">
						<div
							class="bg-primary h-2 rounded-full transition-all"
							style="width: {percentage}%"
						></div>
					</div>
				</div>
			</div>
		{/each}
	{/if}
</div>
