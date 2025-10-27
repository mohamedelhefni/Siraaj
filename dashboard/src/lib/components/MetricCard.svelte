<script>
	import { TrendingUp, TrendingDown, Minus } from 'lucide-svelte';
	import { formatCompactNumber } from '$lib/utils/formatters.js';

	let {
		label,
		currentValue,
		previousValue = null,
		currentPeriod = '',
		previousPeriod = '',
		formatValue = (val) => formatCompactNumber(val),
		isSelected = false,
		isNegativeBetter = false,
		onclick = () => {}
	} = $props();

	// Calculate percentage change
	const change = $derived(() => {
		if (!previousValue || previousValue === 0) return 0;
		return ((currentValue - previousValue) / previousValue) * 100;
	});

	// Get trend icon based on change
	function getTrendIcon(changeVal) {
		if (changeVal > 0) return TrendingUp;
		if (changeVal < 0) return TrendingDown;
		return Minus;
	}

	// Get trend color based on change
	function getTrendColor(changeVal) {
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
>
	<div class="text-muted-foreground mb-1 text-xs font-medium uppercase tracking-wide">
		{label}
	</div>

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
</button>
