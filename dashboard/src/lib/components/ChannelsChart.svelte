<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import * as d3 from 'd3';

	export let data: any[] = [];
	export let chartType: 'bar' | 'pie' = 'pie';

	let chartContainer: HTMLDivElement;
	let tooltip: HTMLDivElement;

	// Channel colors - modern palette
	const channelColors: Record<string, string> = {
		Direct: '#3b82f6', // Blue
		Paid: '#ef4444', // Red
		'google.com': '#4285f4', // Google Blue
		'facebook.com': '#1877f2', // Facebook Blue
		'twitter.com': '#1da1f2', // Twitter Blue
		'linkedin.com': '#0077b5', // LinkedIn Blue
		'youtube.com': '#ff0000', // YouTube Red
		'instagram.com': '#e4405f', // Instagram Pink
		'github.com': '#24292e', // GitHub Dark
		'reddit.com': '#ff4500', // Reddit Orange
		'pinterest.com': '#bd081c', // Pinterest Red
		'tiktok.com': '#000000', // TikTok Black
		Unknown: '#6b7280' // Gray
	};

	// Get color for any domain or fallback to a hash-based color
	function getChannelColor(channel: string): string {
		if (channelColors[channel]) {
			return channelColors[channel];
		}

		// Generate a consistent color based on channel name
		let hash = 0;
		for (let i = 0; i < channel.length; i++) {
			const char = channel.charCodeAt(i);
			hash = (hash << 5) - hash + char;
			hash = hash & hash; // Convert to 32bit integer
		}

		const hue = Math.abs(hash) % 360;
		return `hsl(${hue}, 65%, 50%)`;
	}

	function drawBarChart() {
		if (!chartContainer || data.length === 0) return;

		// Clear previous chart
		d3.select(chartContainer).selectAll('*').remove();

		const margin = { top: 20, right: 20, bottom: 80, left: 100 };
		const width = chartContainer.clientWidth - margin.left - margin.right;
		const height = 400 - margin.top - margin.bottom;

		const svg = d3
			.select(chartContainer)
			.append('svg')
			.attr('width', width + margin.left + margin.right)
			.attr('height', height + margin.top + margin.bottom)
			.style('background', 'transparent');

		const g = svg.append('g').attr('transform', `translate(${margin.left},${margin.top})`);

		// X scale
		const x = d3
			.scaleBand()
			.domain(data.map((d) => d.channel || 'Unknown'))
			.range([0, width])
			.padding(0.3);

		// Y scale
		const maxValue = d3.max(data, (d) => d.total_events || 0) || 0;
		const y = d3
			.scaleLinear()
			.domain([0, maxValue * 1.1])
			.range([height, 0]);

		// Add subtle horizontal grid lines
		g.selectAll('.grid-line')
			.data(y.ticks(5))
			.join('line')
			.attr('class', 'grid-line')
			.attr('x1', 0)
			.attr('x2', width)
			.attr('y1', (d) => y(d))
			.attr('y2', (d) => y(d))
			.attr('stroke', '#f1f5f9')
			.attr('stroke-width', 1)
			.attr('opacity', 0.6);

		// X axis - minimal styling
		g.append('g')
			.attr('transform', `translate(0,${height})`)
			.call(d3.axisBottom(x).tickSize(0))
			.selectAll('text')
			.attr('transform', 'rotate(-45)')
			.style('text-anchor', 'end')
			.style('font-size', '12px')
			.style('font-weight', '500')
			.style('fill', '#64748b');

		// Y axis - minimal styling
		g.append('g')
			.call(
				d3
					.axisLeft(y)
					.ticks(5)
					.tickSize(0)
					.tickFormat((d) => d.toLocaleString())
			)
			.selectAll('text')
			.style('font-size', '12px')
			.style('font-weight', '500')
			.style('fill', '#64748b');

		// Remove axis lines for flat design
		g.selectAll('.domain').remove();

		// Flat bars with no strokes or gradients
		g.selectAll('.bar')
			.data(data)
			.join('rect')
			.attr('class', 'bar')
			.attr('x', (d) => x(d.channel || 'Unknown') || 0)
			.attr('y', height)
			.attr('width', x.bandwidth())
			.attr('height', 0)
			.attr('fill', (d) => getChannelColor(d.channel || 'Unknown'))
			.attr('rx', 4) // Slight rounded corners
			.style('cursor', 'pointer')
			.style('opacity', 0.9)
			.on('mouseenter', function (event, d) {
				d3.select(this)
					.transition()
					.duration(200)
					.style('opacity', 1)
					.attr('transform', 'translateY(-2px)');
				showTooltip(event, d);
			})
			.on('mousemove', (event) => moveTooltip(event))
			.on('mouseleave', function () {
				d3.select(this)
					.transition()
					.duration(200)
					.style('opacity', 0.9)
					.attr('transform', 'translateY(0px)');
				hideTooltip();
			})
			.transition()
			.duration(800)
			.ease(d3.easeBackOut.overshoot(0.1))
			.attr('y', (d) => y(d.total_events || 0))
			.attr('height', (d) => height - y(d.total_events || 0));

		// Value labels
		g.selectAll('.label')
			.data(data)
			.join('text')
			.attr('class', 'label')
			.attr('x', (d) => (x(d.channel || 'Unknown') || 0) + x.bandwidth() / 2)
			.attr('y', height)
			.attr('text-anchor', 'middle')
			.style('font-size', '11px')
			.style('font-weight', '600')
			.style('fill', '#475569')
			.style('opacity', 0)
			.text((d) => (d.total_events || 0).toLocaleString())
			.transition()
			.duration(800)
			.delay(300)
			.ease(d3.easeBackOut)
			.attr('y', (d) => y(d.total_events || 0) - 8)
			.style('opacity', 1);
	}

	function drawPieChart() {
		if (!chartContainer || data.length === 0) return;

		// Clear previous chart
		d3.select(chartContainer).selectAll('*').remove();

		const width = chartContainer.clientWidth;
		const height = 400;
		const radius = Math.min(width, height) / 2 - 80;

		const svg = d3
			.select(chartContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height)
			.style('background', 'transparent');

		const g = svg.append('g').attr('transform', `translate(${width / 2 - 120},${height / 2})`);

		// Calculate total for percentages
		const total = d3.sum(data, (d) => d.total_events || 0);

		// Pie layout with small padding for flat design
		const pie = d3
			.pie<any>()
			.value((d) => d.total_events || 0)
			.sort(null)
			.padAngle(0.005); // Minimal padding for flat look

		// Arc generator - flat donut style
		const arc = d3
			.arc<any>()
			.innerRadius(radius * 0.5) // Thicker donut
			.outerRadius(radius);

		const hoverArc = d3
			.arc<any>()
			.innerRadius(radius * 0.5)
			.outerRadius(radius + 8); // Subtle hover effect

		// Draw slices with flat design
		const slices = g.selectAll('.slice').data(pie(data)).join('g').attr('class', 'slice');

		slices
			.append('path')
			.attr('fill', (d) => getChannelColor(d.data.channel || 'Unknown'))
			.style('cursor', 'pointer')
			.style('opacity', 0.95)
			.on('mouseenter', function (event, d) {
				d3.select(this)
					.transition()
					.duration(200)
					.ease(d3.easeQuadOut)
					.attr('d', hoverArc)
					.style('opacity', 1);
				showTooltip(event, d.data);
			})
			.on('mousemove', (event) => moveTooltip(event))
			.on('mouseleave', function () {
				d3.select(this)
					.transition()
					.duration(200)
					.ease(d3.easeQuadOut)
					.attr('d', arc)
					.style('opacity', 0.95);
				hideTooltip();
			})
			// Smooth pie animation - grow from center
			.attr('d', arc)
			.style('opacity', 0)
			.style('transform', 'scale(0)')
			.transition()
			.duration(800)
			.delay((d, i) => i * 100) // Staggered animation
			.ease(d3.easeBackOut.overshoot(0.3))
			.style('opacity', 0.95)
			.style('transform', 'scale(1)');

		// Add percentage labels with clean styling
		slices
			.append('text')
			.attr('transform', (d) => {
				const [x, y] = arc.centroid(d);
				return `translate(${x},${y})`;
			})
			.attr('text-anchor', 'middle')
			.style('font-size', '12px')
			.style('font-weight', '700')
			.style('fill', '#ffffff')
			.style('text-shadow', '0 1px 2px rgba(0,0,0,0.7)')
			.style('pointer-events', 'none')
			.text((d) => {
				const percentage = ((d.data.total_events / total) * 100).toFixed(1);
				return percentage > 4 ? `${percentage}%` : '';
			})
			.style('opacity', 0)
			.transition()
			.duration(600)
			.delay((d, i) => i * 100 + 400)
			.ease(d3.easeQuadOut)
			.style('opacity', 1);

		// Center label showing total with clean design
		const centerGroup = g.append('g').attr('class', 'center-label');

		centerGroup
			.append('text')
			.attr('text-anchor', 'middle')
			.attr('dy', '-0.3em')
			.style('font-size', '28px')
			.style('font-weight', '800')
			.style('fill', '#1e293b')
			.text(total.toLocaleString())
			.style('opacity', 0)
			.transition()
			.duration(600)
			.delay(800)
			.ease(d3.easeQuadOut)
			.style('opacity', 1);

		centerGroup
			.append('text')
			.attr('text-anchor', 'middle')
			.attr('dy', '1.2em')
			.style('font-size', '14px')
			.style('font-weight', '500')
			.style('fill', '#64748b')
			.text('Total Events')
			.style('opacity', 0)
			.transition()
			.duration(600)
			.delay(900)
			.ease(d3.easeQuadOut)
			.style('opacity', 1);

		// Clean legend design
		const legend = svg.append('g').attr('transform', `translate(${width - 200}, 40)`);

		const legendItems = legend
			.selectAll('.legend-item')
			.data(data.slice(0, 8)) // Limit to top 8 for space
			.join('g')
			.attr('class', 'legend-item')
			.attr('transform', (d, i) => `translate(0, ${i * 24})`)
			.style('cursor', 'pointer')
			.style('opacity', 0);

		legendItems
			.append('rect')
			.attr('x', 0)
			.attr('y', 0)
			.attr('width', 12)
			.attr('height', 12)
			.attr('rx', 2)
			.attr('fill', (d) => getChannelColor(d.channel || 'Unknown'));

		legendItems
			.append('text')
			.attr('x', 18)
			.attr('y', 6)
			.attr('dy', '0.35em')
			.style('font-size', '12px')
			.style('font-weight', '500')
			.style('fill', '#1e293b')
			.text((d) => {
				const percentage = ((d.total_events / total) * 100).toFixed(1);
				const label = d.channel || 'Unknown';
				return `${label.length > 15 ? label.substring(0, 15) + '...' : label} (${percentage}%)`;
			});

		// Animate legend items
		legendItems
			.transition()
			.duration(400)
			.delay((d, i) => i * 80 + 600)
			.ease(d3.easeQuadOut)
			.style('opacity', 1);

		// Add legend interactions
		legendItems
			.on('mouseenter', function (event, d) {
				// Highlight corresponding pie slice
				const sliceIndex = data.findIndex((item) => item.channel === d.channel);
				if (sliceIndex >= 0) {
					g.selectAll('.slice path').style('opacity', 0.4);
					g.select(`.slice:nth-child(${sliceIndex + 1}) path`)
						.style('opacity', 1)
						.transition()
						.duration(150)
						.attr('d', hoverArc);
				}
				showTooltip(event, d);
			})
			.on('mouseleave', function () {
				g.selectAll('.slice path').style('opacity', 0.95).transition().duration(150).attr('d', arc);
				hideTooltip();
			});
	}

	function showTooltip(event: MouseEvent, d: any) {
		const total = d3.sum(data, (item) => item.total_events || 0);
		const percentage = ((d.total_events / total) * 100).toFixed(1);

		tooltip.innerHTML = `
			<div class="tooltip-header">
				<div class="tooltip-indicator" style="background-color: ${getChannelColor(d.channel || 'Unknown')}"></div>
				<span class="tooltip-title">${d.channel || 'Unknown'}</span>
				<span class="tooltip-percentage">${percentage}%</span>
			</div>
			<div class="tooltip-metrics">
				<div class="tooltip-metric">
					<span class="metric-label">Events</span>
					<span class="metric-value">${(d.total_events || 0).toLocaleString()}</span>
				</div>
				<div class="tooltip-metric">
					<span class="metric-label">Users</span>
					<span class="metric-value">${(d.unique_users || 0).toLocaleString()}</span>
				</div>
				<div class="tooltip-metric">
					<span class="metric-label">Visits</span>
					<span class="metric-value">${(d.total_visits || 0).toLocaleString()}</span>
				</div>
				<div class="tooltip-metric">
					<span class="metric-label">Page Views</span>
					<span class="metric-value">${(d.page_views || 0).toLocaleString()}</span>
				</div>
			</div>
		`;
		tooltip.style.display = 'block';
		tooltip.style.opacity = '0';
		moveTooltip(event);

		// Fade in animation
		requestAnimationFrame(() => {
			tooltip.style.opacity = '1';
		});
	}

	function moveTooltip(event: MouseEvent) {
		tooltip.style.left = event.pageX + 10 + 'px';
		tooltip.style.top = event.pageY + 10 + 'px';
	}

	function hideTooltip() {
		tooltip.style.display = 'none';
	}

	function renderChart() {
		if (chartType === 'bar') {
			drawBarChart();
		} else {
			drawPieChart();
		}
	}

	// Redraw chart when data or chartType changes
	let lastDataString = '';
	let lastChartType = '';

	$: {
		const currentDataString = JSON.stringify(data);
		if (chartContainer && (currentDataString !== lastDataString || chartType !== lastChartType)) {
			lastDataString = currentDataString;
			lastChartType = chartType;
			renderChart();
		}
	}

	onMount(() => {
		if (data.length > 0) {
			renderChart();
		}

		// Redraw on window resize
		const handleResize = () => {
			if (data.length > 0) {
				renderChart();
			}
		};
		window.addEventListener('resize', handleResize);

		return () => {
			window.removeEventListener('resize', handleResize);
		};
	});

	// Remove afterUpdate to prevent unnecessary re-renders
</script>

<div class="chart-container" bind:this={chartContainer}></div>

<!-- Tooltip -->
<div bind:this={tooltip} class="tooltip"></div>

<style>
	.chart-container {
		min-height: 400px;
		width: 100%;
		position: relative;
	}

	:global(.tooltip) {
		position: absolute;
		display: none;
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0.9));
		backdrop-filter: blur(10px);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 12px;
		padding: 0;
		font-size: 13px;
		pointer-events: none;
		z-index: 1000;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
		transition: opacity 0.2s ease;
		max-width: 250px;
	}

	:global(.tooltip-header) {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 16px 8px 16px;
		border-bottom: 1px solid rgba(0, 0, 0, 0.06);
		margin-bottom: 8px;
	}

	:global(.tooltip-indicator) {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		flex-shrink: 0;
	}

	:global(.tooltip-title) {
		font-weight: 600;
		color: #1e293b;
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	:global(.tooltip-percentage) {
		font-weight: 700;
		color: #3b82f6;
		font-size: 12px;
	}

	:global(.tooltip-metrics) {
		padding: 0 16px 12px 16px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	:global(.tooltip-metric) {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	:global(.metric-label) {
		color: #64748b;
		font-size: 12px;
		font-weight: 500;
	}

	:global(.metric-value) {
		color: #1e293b;
		font-weight: 600;
		font-size: 12px;
	}

	:global(.bar) {
		cursor: pointer;
		transition: all 0.2s ease;
	}

	:global(.slice path) {
		cursor: pointer;
		transition: all 0.3s ease;
	}

	:global(.legend-item) {
		transition: all 0.2s ease;
	}

	:global(.legend-item:hover) {
		opacity: 0.8;
	}
</style>
