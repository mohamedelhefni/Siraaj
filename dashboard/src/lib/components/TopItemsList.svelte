<script>
	import { getCountryFlag, getBrowserName, getFaviconUrl, getSourceDisplayName } from '$lib/utils/icons';
	import { Chrome, Globe } from 'lucide-svelte';

	let { items = [], labelKey = 'name', valueKey = 'count', maxItems = 10, type = 'default' } = $props();

	const displayItems = $derived(items.slice(0, maxItems));
	const total = $derived(items.reduce((sum, item) => sum + (item[valueKey] || 0), 0));

	// Browser icon map
	const browserIcons = {
		chrome: '🌐',
		firefox: '🦊',
		safari: '🧭',
		edge: '🌊',
		opera: '⭕',
		brave: '🦁',
		unknown: '🌍'
	};
</script>

<div class="space-y-2">
	{#if displayItems.length === 0}
		<p class="text-sm text-muted-foreground text-center py-4">No data available</p>
	{:else}
		{#each displayItems as item, i}
			{@const percentage = total > 0 ? ((item[valueKey] / total) * 100).toFixed(1) : 0}
			{@const displayLabel = type === 'source' ? getSourceDisplayName(item[labelKey]) : item[labelKey]}
			{@const faviconUrl = type === 'source' ? getFaviconUrl(item[labelKey]) : null}
			{@const countryFlag = type === 'country' ? getCountryFlag(item[labelKey]) : null}
			{@const browserIcon = type === 'browser' ? browserIcons[getBrowserName(item[labelKey])] : null}

			<div class="flex items-center justify-between space-x-2">
				<div class="flex-1 min-w-0">
					<div class="flex items-center justify-between mb-1">
						<div class="flex items-center gap-2 min-w-0 flex-1">
							{#if countryFlag}
								<span class="text-lg flex-shrink-0">{countryFlag}</span>
							{:else if browserIcon}
								<span class="text-base flex-shrink-0">{browserIcon}</span>
							{:else if faviconUrl}
								<img src={faviconUrl} alt="" class="w-4 h-4 flex-shrink-0" />
							{:else if type === 'source' && item[labelKey] === 'Direct'}
								<span class="text-base flex-shrink-0">🔗</span>
							{/if}
							<span class="text-sm font-medium truncate" title={displayLabel}>
								{i + 1}. {displayLabel}
							</span>
						</div>
						<span class="text-sm text-muted-foreground ml-2 flex-shrink-0">
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
