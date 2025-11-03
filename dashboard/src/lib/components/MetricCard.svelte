<script lang="ts">
	import { TrendingUp, TrendingDown, Minus } from 'lucide-svelte';
	import { formatCompactNumber } from '$lib/utils/formatters.js';

	let {
		label,
		currentValue,
		previousValue = null,
		currentPeriod = '',
		previousPeriod = '',
		formatValue = (val: number) => formatCompactNumber(val),
		isSelected = false,
		isNegativeBetter = false,
		loading = false,
		onclick = () => {}
	}: {
		label: string;
		currentValue: number;
		previousValue?: number | null;
		currentPeriod?: string;
		previousPeriod?: string;
		formatValue?: (val: number) => string;
		isSelected?: boolean;
		isNegativeBetter?: boolean;
		loading?: boolean;
		onclick?: () => void;
	} = $props();

	// Calculate percentage change
	const change = $derived(() => {
		if (!previousValue || previousValue === 0) return 0;
		return ((currentValue - previousValue) / previousValue) * 100;
	});

	// Get trend icon based on change
	function getTrendIcon(changeVal: number) {
		if (changeVal > 0) return TrendingUp;
		if (changeVal < 0) return TrendingDown;
		return Minus;
	}

	// Get trend color based on change
	function getTrendColor(changeVal: number) {
		if (changeVal === 0) return 'text-gray-600';

		if (isNegativeBetter) {
			// For metrics like bounce rate, lower is better
			return changeVal < 0 ? 'text-green-600' : 'text-red-600';
		} else {
			// For most metrics, higher is better
			return changeVal > 0 ? 'text-green-600' : 'text-red-600';
		}
	}
</script>

<button
	class="hover:bg-accent group border-b border-r p-4 text-left transition-colors lg:border-b-0 {isSelected
		? 'bg-accent'
		: ''}"
	{onclick}
	disabled={loading}
>
	<div class="text-muted-foreground mb-1 text-xs font-medium uppercase tracking-wide">
		{label}
	</div>

	{#if loading}
		<!-- Loading Skeleton -->
		<div class="mb-1">
			<div class="bg-muted h-8 w-24 animate-pulse rounded"></div>
			{#if currentPeriod}
				<div class="bg-muted mt-1 h-3 w-16 animate-pulse rounded"></div>
			{/if}
		</div>

		{#if previousValue !== null}
			<div class="mt-2 opacity-60">
				<div class="bg-muted h-6 w-20 animate-pulse rounded"></div>
				{#if previousPeriod}
					<div class="bg-muted mt-1 h-3 w-16 animate-pulse rounded"></div>
				{/if}
			</div>
		{/if}
	{:else}
		<!-- Current Period -->
		<div class="mb-1">
			<div class="text-2xl font-bold">{formatValue(currentValue)}</div>
			{#if currentPeriod}
				<div class="text-muted-foreground text-xs">{currentPeriod}</div>
			{/if}
		</div>

		<!-- Previous Period (with fade) -->
		{#if previousValue !== null}
			<div class="mt-2 opacity-60">
				<div class="text-lg font-semibold">{formatValue(previousValue)}</div>
				{#if previousPeriod}
					<div class="text-muted-foreground text-xs">{previousPeriod}</div>
				{/if}
			</div>

			<!-- Change Indicator -->
			{#if change() !== 0}
				{@const TrendIcon = getTrendIcon(change())}
				<div class="text-xs {getTrendColor(change())} mt-1 flex items-center gap-1">
					<TrendIcon class="h-3 w-3" />
					{Math.abs(change()).toFixed(0)}%
				</div>
			{/if}
		{/if}
	{/if}
</button>
