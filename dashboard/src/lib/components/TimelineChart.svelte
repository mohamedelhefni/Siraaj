<script>
	import { onMount } from 'svelte';
	import * as d3 from 'd3';
	import { Button } from '$lib/components/ui/button';
	import { Eye, EyeOff } from 'lucide-svelte';

	let { data = [], comparisonData = [], format = 'day', metric = 'events' } = $props();
	let svgElement = $state();
	let chartContainer = $state();
	let showComparison = $state(true);

	// Get label based on metric type
	const metricLabels = {
		events: 'Events',
		users: 'Unique Visitors',
		visits: 'Total Visits',
		page_views: 'Page Views',
		views_per_visit: 'Views per Visit',
		bounce_rate: 'Bounce Rate',
		visit_duration: 'Visit Duration'
	};

	// Format the count value based on metric type
	function formatCount(count, metricType) {
		if (metricType === 'bounce_rate') {
			return count.toFixed(1) + '%';
		} else if (metricType === 'visit_duration') {
			// Format seconds into readable duration
			if (count < 60) {
				return Math.floor(count) + 's';
			} else if (count < 3600) {
				const minutes = Math.floor(count / 60);
				const seconds = Math.floor(count % 60);
				return `${minutes}m ${seconds}s`;
			} else {
				const hours = Math.floor(count / 3600);
				const minutes = Math.floor((count % 3600) / 60);
				return `${hours}h ${minutes}m`;
			}
		} else if (metricType === 'views_per_visit') {
			return count.toFixed(2);
		} else {
			return Math.round(count).toLocaleString();
		}
	}

	// Format date labels based on granularity
	function formatLabel(dateStr) {
		const date = new Date(dateStr);

		if (format === 'hour') {
			return d3.timeFormat('%H:%M')(date);
		} else if (format === 'month') {
			return d3.timeFormat('%b %Y')(date);
		} else {
			return d3.timeFormat('%b %d')(date);
		}
	}

	function formatTooltipDate(dateStr) {
		const date = new Date(dateStr);

		if (format === 'hour') {
			return d3.timeFormat('%b %d, %H:%M')(date);
		} else if (format === 'month') {
			return d3.timeFormat('%B %Y')(date);
		} else {
			return d3.timeFormat('%B %d, %Y')(date);
		}
	}

	function drawChart() {
		if (!svgElement || !chartContainer || data.length === 0) return;

		// Clear previous chart
		d3.select(svgElement).selectAll('*').remove();

		const margin = { top: 20, right: 30, bottom: 30, left: 50 };
		const width = chartContainer.clientWidth - margin.left - margin.right;
		const height = 300 - margin.top - margin.bottom;

		const svg = d3
			.select(svgElement)
			.attr('width', width + margin.left + margin.right)
			.attr('height', height + margin.top + margin.bottom)
			.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		// Parse dates and prepare data
		const parsedData = data.map((d) => ({
			date: new Date(d.date),
			count: d.count
		}));

		// Parse comparison data if available
		const parsedComparisonData = comparisonData.map((d) => ({
			date: new Date(d.date),
			count: d.count
		}));

		// Create scales
		const x = d3
			.scaleTime()
			.domain(d3.extent(parsedData, (d) => d.date))
			.range([0, width]);

		// Calculate max value from both datasets for consistent y-axis
		const maxValue = d3.max([
			d3.max(parsedData, (d) => d.count),
			showComparison && parsedComparisonData.length > 0
				? d3.max(parsedComparisonData, (d) => d.count)
				: 0
		]);

		const y = d3.scaleLinear().domain([0, maxValue]).nice().range([height, 0]);

		// Create line generator
		const line = d3
			.line()
			.x((d) => x(d.date))
			.y((d) => y(d.count))
			.curve(d3.curveMonotoneX);

		// Create area generator for fill
		const area = d3
			.area()
			.x((d) => x(d.date))
			.y0(height)
			.y1((d) => y(d.count))
			.curve(d3.curveMonotoneX);

		// Add X axis
		svg
			.append('g')
			.attr('transform', `translate(0,${height})`)
			.call(
				d3
					.axisBottom(x)
					.ticks(Math.min(parsedData.length, 7))
					.tickFormat((d) => formatLabel(d))
			)
			.selectAll('text')
			.style('font-size', '12px')
			.style('fill', '#6b7280');

		// Add Y axis
		const yAxis = d3.axisLeft(y).ticks(5);

		// Format Y axis based on metric type
		if (metric === 'bounce_rate') {
			yAxis.tickFormat((d) => d.toFixed(0) + '%');
		} else if (metric === 'visit_duration') {
			yAxis.tickFormat((d) => {
				if (d < 60) return Math.floor(d) + 's';
				const minutes = Math.floor(d / 60);
				return minutes + 'm';
			});
		} else if (metric === 'views_per_visit') {
			yAxis.tickFormat((d) => d.toFixed(1));
		} else {
			yAxis.tickFormat((d) => d.toLocaleString());
		}

		svg
			.append('g')
			.call(yAxis)
			.selectAll('text')
			.style('font-size', '12px')
			.style('fill', '#6b7280');

		// Add grid lines
		svg
			.append('g')
			.attr('class', 'grid')
			.attr('opacity', 0.1)
			.call(d3.axisLeft(y).ticks(5).tickSize(-width).tickFormat(''));

		// Define linear gradient for area fill
		const gradient = svg
			.append('defs')
			.append('linearGradient')
			.attr('id', 'area-gradient')
			.attr('x1', '0%')
			.attr('x2', '0%')
			.attr('y1', '0%')
			.attr('y2', '100%');

		gradient
			.append('stop')
			.attr('offset', '0%')
			.attr('stop-color', 'rgb(99, 102, 241)')
			.attr('stop-opacity', 0.1);

		gradient
			.append('stop')
			.attr('offset', '100%')
			.attr('stop-color', 'rgb(99, 102, 241)')
			.attr('stop-opacity', 0);

		// Add area with animation and gradient
		svg
			.append('path')
			.datum(parsedData)
			.attr('fill', 'url(#area-gradient)')
			.attr('d', area)
			.style('opacity', 0)
			.transition()
			.duration(750)
			.ease(d3.easeQuadInOut)
			.style('opacity', 1);

		// Add line with animation
		const path = svg
			.append('path')
			.datum(parsedData)
			.attr('fill', 'none')
			.attr('stroke', 'rgb(99, 102, 241)')
			.attr('stroke-width', 2)
			.attr('d', line);

		const totalLength = path.node().getTotalLength();

		path
			.attr('stroke-dasharray', totalLength + ' ' + totalLength)
			.attr('stroke-dashoffset', totalLength)
			.transition()
			.duration(750)
			.ease(d3.easeQuadInOut)
			.attr('stroke-dashoffset', 0);

		// Add comparison line (dotted, previous period)
		if (showComparison && parsedComparisonData.length > 0) {
			// Shift comparison dates to align with current period
			const timeDiff = parsedData[0].date.getTime() - parsedComparisonData[0].date.getTime();
			const alignedComparisonData = parsedComparisonData.map((d) => ({
				date: new Date(d.date.getTime() + timeDiff),
				count: d.count,
				originalDate: d.date
			}));

			const comparisonLine = d3
				.line()
				.x((d) => x(d.date))
				.y((d) => y(d.count))
				.curve(d3.curveMonotoneX);

			const comparisonPath = svg
				.append('path')
				.datum(alignedComparisonData)
				.attr('fill', 'none')
				.attr('stroke', 'rgb(156, 163, 175)') // gray-400
				.attr('stroke-width', 2)
				.attr('stroke-dasharray', '5,5') // dotted line
				.attr('opacity', 0.6)
				.attr('d', comparisonLine);

			const comparisonLength = comparisonPath.node().getTotalLength();

			comparisonPath
				.attr('stroke-dasharray', comparisonLength + ' ' + comparisonLength)
				.attr('stroke-dashoffset', comparisonLength)
				.transition()
				.duration(750)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0)
				.attr('stroke-dasharray', '5,5'); // Reset to dotted after animation

			// Add comparison dots
			const comparisonDots = svg
				.selectAll('.comparison-dot')
				.data(alignedComparisonData)
				.enter()
				.append('circle')
				.attr('class', 'comparison-dot')
				.attr('cx', (d) => x(d.date))
				.attr('cy', (d) => y(d.count))
				.attr('r', 0)
				.attr('fill', 'rgb(156, 163, 175)')
				.attr('opacity', 0.6)
				.style('cursor', 'pointer');

			comparisonDots
				.transition()
				.delay((d, i) => i * 50)
				.duration(300)
				.attr('r', 3);

			// Add hover effects f0or comparison dots
			comparisonDots
				.on('mouseover', function (event, d) {
					d3.select(this).transition().duration(200).attr('r', 5);

					tooltip.transition().duration(200).style('opacity', 1);

					tooltip
						.html(
							`<strong>Previous Period</strong><br/>${formatTooltipDate(d.originalDate.toISOString())}<br/>${metricLabels[metric] || 'Count'}: <strong>${formatCount(d.count, metric)}</strong>`
						)
						.style('left', event.pageX + 10 + 'px')
						.style('top', event.pageY - 10 + 'px');
				})
				.on('mouseout', function () {
					d3.select(this).transition().duration(200).attr('r', 3);

					tooltip.transition().duration(200).style('opacity', 0);
				});
		}

		// Add dots
		const dots = svg
			.selectAll('.dot')
			.data(parsedData)
			.enter()
			.append('circle')
			.attr('class', 'dot')
			.attr('cx', (d) => x(d.date))
			.attr('cy', (d) => y(d.count))
			.attr('r', 0)
			.attr('fill', 'rgb(99, 102, 241)')
			.style('cursor', 'pointer');

		dots
			.transition()
			.delay((d, i) => i * 50)
			.duration(300)
			.attr('r', 4);

		// Create tooltip
		const tooltip = d3
			.select('body')
			.append('div')
			.style('position', 'absolute')
			.style('background', 'rgba(0, 0, 0, 0.8)')
			.style('color', 'white')
			.style('padding', '8px 12px')
			.style('border-radius', '4px')
			.style('font-size', '14px')
			.style('pointer-events', 'none')
			.style('opacity', 0)
			.style('z-index', '1000');

		// Add hover effects
		dots
			.on('mouseover', function (event, d) {
				d3.select(this).transition().duration(200).attr('r', 6);

				tooltip.transition().duration(200).style('opacity', 1);

				tooltip
					.html(
						`<strong>${formatTooltipDate(d.date.toISOString())}</strong><br/>${metricLabels[metric] || 'Count'}: <strong>${formatCount(d.count, metric)}</strong>`
					)
					.style('left', event.pageX + 10 + 'px')
					.style('top', event.pageY - 10 + 'px');
			})
			.on('mouseout', function () {
				d3.select(this).transition().duration(200).attr('r', 4);

				tooltip.transition().duration(200).style('opacity', 0);
			});

		// Cleanup tooltip on destroy
		return () => {
			tooltip.remove();
		};
	}

	onMount(() => {
		const cleanup = drawChart();

		// Handle window resize
		const handleResize = () => {
			drawChart();
		};

		window.addEventListener('resize', handleResize);

		return () => {
			window.removeEventListener('resize', handleResize);
			if (cleanup) cleanup();
		};
	});

	$effect(() => {
		// Redraw when data or metric changes
		if (data && metric) {
			drawChart();
		}
	});

	$effect(() => {
		// Redraw when showComparison changes
		if (showComparison !== undefined) {
			drawChart();
		}
	});
</script>

<div class="relative">
	{#if data.length > 0 && comparisonData.length > 0}
		<div class="absolute -bottom-5 right-0 z-10 flex items-center gap-4">
			<div class="flex items-center gap-3 text-sm">
				<div class="flex items-center gap-2">
					<div class="h-0.5 w-6 bg-indigo-500"></div>
					<span class="text-muted-foreground">Current Period</span>
				</div>
				<div class="flex items-center gap-2">
					<div class="h-0.5 w-6 border-t-2 border-dashed border-gray-400"></div>
					<span class="text-muted-foreground">Previous Period</span>
				</div>
			</div>
			<Button
				variant="ghost"
				size="sm"
				onclick={() => (showComparison = !showComparison)}
				class="gap-2"
			>
				{#if showComparison}
					<Eye class="h-4 w-4" />
					Hide Comparison
				{:else}
					<EyeOff class="h-4 w-4" />
					Show Comparison
				{/if}
			</Button>
		</div>
	{/if}

	<div bind:this={chartContainer} class="h-auto w-full">
		{#if data.length === 0}
			<div class="text-muted-foreground flex h-full items-center justify-center">
				<p>No data available for the selected period</p>
			</div>
		{:else}
			<svg bind:this={svgElement} class="w-full"></svg>
		{/if}
	</div>
</div>
