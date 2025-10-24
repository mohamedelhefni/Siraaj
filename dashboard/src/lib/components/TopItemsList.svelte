<script>
	let { items = [], labelKey = 'name', valueKey = 'count', maxItems = 10 } = $props();

	const displayItems = $derived(items.slice(0, maxItems));
	const total = $derived(items.reduce((sum, item) => sum + (item[valueKey] || 0), 0));
</script>

<div class="space-y-2">
	{#if displayItems.length === 0}
		<p class="text-muted-foreground py-4 text-center text-sm">No data available</p>
	{:else}
		{#each displayItems as item, i}
			{@const percentage = total > 0 ? ((item[valueKey] / total) * 100).toFixed(1) : 0}
			<div class="flex items-center justify-between space-x-2">
				<div class="min-w-0 flex-1">
					<div class="mb-1 flex items-center justify-between">
						<span class="truncate text-sm font-medium" title={item[labelKey]}>
							{i + 1}. {item[labelKey]}
						</span>
						<span class="text-muted-foreground ml-2 text-sm">
							{item[valueKey]?.toLocaleString()} ({percentage}%)
						</span>
					</div>
					<div class="bg-secondary h-2 w-full rounded-full">
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
