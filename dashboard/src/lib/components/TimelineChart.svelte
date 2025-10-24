<script>
	import { onMount } from 'svelte';
	import * as d3 from 'd3';

	let { data = [], format = 'day', metric = 'events' } = $props();
	let svgElement = $state();
	let chartContainer = $state();

	// Get label based on metric type
	const metricLabels = {
		events: 'Events',
		users: 'Unique Visitors',
		visits: 'Total Visits',
		page_views: 'Page Views'
	};

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

		// Create scales
		const x = d3
			.scaleTime()
			.domain(d3.extent(parsedData, (d) => d.date))
			.range([0, width]);

		const y = d3
			.scaleLinear()
			.domain([0, d3.max(parsedData, (d) => d.count)])
			.nice()
			.range([height, 0]);

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
		svg
			.append('g')
			.call(d3.axisLeft(y).ticks(5))
			.selectAll('text')
			.style('font-size', '12px')
			.style('fill', '#6b7280');

		// Add grid lines
		svg
			.append('g')
			.attr('class', 'grid')
			.attr('opacity', 0.1)
			.call(d3.axisLeft(y).ticks(5).tickSize(-width).tickFormat(''));

		// Add area with animation
		svg
			.append('path')
			.datum(parsedData)
			.attr('fill', 'rgba(99, 102, 241, 0.1)')
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
						`<strong>${formatTooltipDate(d.date.toISOString())}</strong><br/>${metricLabels[metric]}: <strong>${d.count}</strong>`
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
</script>

<div bind:this={chartContainer} class="h-[300px] w-full">
	{#if data.length === 0}
		<div class="text-muted-foreground flex h-full items-center justify-center">
			<p>No data available for the selected period</p>
		</div>
	{:else}
		<svg bind:this={svgElement} class="w-full"></svg>
	{/if}
</div>
