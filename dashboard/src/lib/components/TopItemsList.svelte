<script>
	import {
		getCountryFlag,
		getBrowserName,
		getFaviconUrl,
		getSourceDisplayName
	} from '$lib/utils/icons';
	import { Chrome, Globe } from 'lucide-svelte';

	let {
		items = [],
		labelKey = 'name',
		valueKey = 'count',
		maxItems = 10,
		type = 'default',
		onclick = null
	} = $props();

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

	function handleClick(item) {
		if (onclick) {
			onclick(item);
		}
	}
</script>

<div class="space-y-2">
	{#if displayItems.length === 0}
		<p class="text-muted-foreground py-4 text-center text-sm">No data available</p>
	{:else}
		{#each displayItems as item, i}
			{@const percentage = total > 0 ? ((item[valueKey] / total) * 100).toFixed(1) : 0}
			{@const displayLabel =
				type === 'source' ? getSourceDisplayName(item[labelKey]) : item[labelKey]}
			{@const faviconUrl = type === 'source' ? getFaviconUrl(item[labelKey]) : null}
			{@const countryFlag = type === 'country' ? getCountryFlag(item[labelKey]) : null}
			{@const browserIcon =
				type === 'browser' ? browserIcons[getBrowserName(item[labelKey])] : null}

			<div
				class="flex items-center justify-between space-x-2 {onclick
					? 'hover:bg-accent -m-1 cursor-pointer rounded-md p-1 transition-colors'
					: ''}"
				onclick={() => handleClick(item)}
				onkeydown={(e) => e.key === 'Enter' && handleClick(item)}
				role="button"
				tabindex="0"
			>
				<div class="min-w-0 flex-1">
					<div class="mb-1 flex items-center justify-between">
						<div class="flex min-w-0 flex-1 items-center gap-2">
							{#if countryFlag}
								<span class="shrink-0 text-lg">{countryFlag}</span>
							{:else if browserIcon}
								<span class="shrink-0 text-base">{browserIcon}</span>
							{:else if faviconUrl}
								<img src={faviconUrl} alt="" class="h-4 w-4 shrink-0" />
							{:else if type === 'source' && item[labelKey] === 'Direct'}
								<span class="shrink-0 text-base">🔗</span>
							{/if}
							<span class="truncate text-sm font-medium" title={displayLabel}>
								{i + 1}. {displayLabel}
							</span>
						</div>
						<span class="text-muted-foreground ml-2 shrink-0 text-sm">
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
