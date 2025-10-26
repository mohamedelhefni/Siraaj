<script>
	import {
		getCountryFlag,
		getBrowserName,
		getFaviconUrl,
		getSourceDisplayName
	} from '$lib/utils/icons';
	import { Chrome, Globe, ChevronDown } from 'lucide-svelte';
	import Modal from './Modal.svelte';

	let {
		items = [],
		labelKey = 'name',
		valueKey = 'count',
		maxItems = 10,
		type = 'default',
		onclick = null,
		showMoreTitle = 'All Items'
	} = $props();

	let modalOpen = $state(false);

	const displayItems = $derived(items.slice(0, maxItems));
	const hasMore = $derived(items.length > maxItems);
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

	// Device icon map
	const deviceIcons = {
		desktop: '🖥️',
		mobile: '📱',
		tablet: '📱',
		unknown: '💻'
	};

	// OS icon map
	const osIcons = {
		windows: '🪟',
		macos: '🍎',
		linux: '🐧',
		android: '🤖',
		ios: '📱',
		unknown: '💻'
	};

	function getDeviceIcon(device) {
		const deviceLower = (device || '').toLowerCase();
		if (deviceLower.includes('desktop')) return deviceIcons.desktop;
		if (deviceLower.includes('mobile') || deviceLower.includes('phone')) return deviceIcons.mobile;
		if (deviceLower.includes('tablet')) return deviceIcons.tablet;
		return deviceIcons.unknown;
	}

	function getOSIcon(os) {
		const osLower = (os || '').toLowerCase();
		if (osLower.includes('windows')) return osIcons.windows;
		if (osLower.includes('mac') || osLower.includes('darwin')) return osIcons.macos;
		if (osLower.includes('linux')) return osIcons.linux;
		if (osLower.includes('android')) return osIcons.android;
		if (osLower.includes('ios') || osLower.includes('iphone') || osLower.includes('ipad'))
			return osIcons.ios;
		return osIcons.unknown;
	}

	function handleClick(item) {
		if (onclick) {
			onclick(item);
		}
	}

	function openModal() {
		modalOpen = true;
	}

	function renderItem(item, index) {
		const percentage = total > 0 ? ((item[valueKey] / total) * 100).toFixed(1) : 0;
		const displayLabel = type === 'source' ? getSourceDisplayName(item[labelKey]) : item[labelKey];
		const faviconUrl = type === 'source' ? getFaviconUrl(item[labelKey]) : null;
		const countryFlag = type === 'country' ? getCountryFlag(item[labelKey]) : null;
		const browserIcon = type === 'browser' ? browserIcons[getBrowserName(item[labelKey])] : null;
		const deviceIcon = type === 'device' ? getDeviceIcon(item[labelKey]) : null;
		const osIcon = type === 'os' ? getOSIcon(item[labelKey]) : null;

		return {
			item,
			index,
			percentage,
			displayLabel,
			faviconUrl,
			countryFlag,
			browserIcon,
			deviceIcon,
			osIcon
		};
	}
</script>

<div class="space-y-2">
	{#if displayItems.length === 0}
		<p class="text-muted-foreground py-4 text-center text-sm">No data available</p>
	{:else}
		{#each displayItems as item, i}
			{@const {
				percentage,
				displayLabel,
				faviconUrl,
				countryFlag,
				browserIcon,
				deviceIcon,
				osIcon
			} = renderItem(item, i)}

			<div
				class="my-1 flex items-center justify-between space-x-2 p-2 {onclick
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
							{:else if deviceIcon}
								<span class="shrink-0 text-base">{deviceIcon}</span>
							{:else if osIcon}
								<span class="shrink-0 text-base">{osIcon}</span>
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

		{#if hasMore}
			<button
				type="button"
				onclick={openModal}
				class="text-muted-foreground hover:bg-accent hover:text-foreground mt-2 flex w-full items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors"
			>
				<ChevronDown class="h-4 w-4" />
				Show More ({items.length - maxItems} more)
			</button>
		{/if}
	{/if}
</div>

<Modal bind:open={modalOpen} title={showMoreTitle}>
	{#snippet children()}
		<div class="space-y-2">
			{#each items as item, i}
				{@const {
					percentage,
					displayLabel,
					faviconUrl,
					countryFlag,
					browserIcon,
					deviceIcon,
					osIcon
				} = renderItem(item, i)}

				<div
					class="my-1 flex items-center justify-between space-x-2 p-2 {onclick
						? 'hover:bg-accent -m-1 cursor-pointer rounded-md p-1 transition-colors'
						: ''}"
					onclick={() => {
						handleClick(item);
						modalOpen = false;
					}}
					onkeydown={(e) => {
						if (e.key === 'Enter') {
							handleClick(item);
							modalOpen = false;
						}
					}}
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
								{:else if deviceIcon}
									<span class="shrink-0 text-base">{deviceIcon}</span>
								{:else if osIcon}
									<span class="shrink-0 text-base">{osIcon}</span>
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
		</div>
	{/snippet}
</Modal>
