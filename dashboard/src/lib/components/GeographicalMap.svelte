<script>
	import { onMount } from 'svelte';
	import * as d3 from 'd3';
	import * as topojson from 'topojson-client';

	let { data = [], onclick = null } = $props();
	let svgElement = $state();
	let mapContainer = $state();

	// World map TopoJSON URL (using Natural Earth data)
	const WORLD_MAP_URL = 'https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json';

	// Helper function to normalize country names
	function normalizeCountryName(name) {
		if (!name) return name;
		// Replace Israel with Palestine
		if (name.toLowerCase() === 'israel') {
			return 'Palestine';
		}
		return name;
	}

	async function drawMap() {
		if (!svgElement || !mapContainer) return;

		// Clear previous map
		d3.select(svgElement).selectAll('*').remove();

		const width = mapContainer.clientWidth;
		const height = 400;

		const svg = d3
			.select(svgElement)
			.attr('width', width)
			.attr('height', height)
			.attr('viewBox', [0, 0, width, height]);

		// Create projection
		const projection = d3
			.geoNaturalEarth1()
			.scale(width / 6)
			.translate([width / 2, height / 2]);

		const path = d3.geoPath(projection);

		// Create color scale based on data
		const maxCount = d3.max(data, (d) => d.count) || 1;
		const colorScale = d3.scaleSequential(d3.interpolateBlues).domain([0, maxCount]);

		// Create a map of country data for quick lookup
		// Normalize country names (replace Israel with Palestine)
		const countryDataMap = new Map(
			data.map((d) => [normalizeCountryName(d.country).toLowerCase(), d.count])
		);

		try {
			// Load world map data
			const worldData = await d3.json(WORLD_MAP_URL);
			const countries = topojson.feature(worldData, worldData.objects.countries);

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

			// Draw countries
			svg
				.append('g')
				.selectAll('path')
				.data(countries.features)
				.join('path')
				.attr('d', path)
				.attr('fill', (d) => {
					// Normalize country name (replace Israel with Palestine)
					let countryName = normalizeCountryName(d.properties.name);
					const count = countryDataMap.get(countryName?.toLowerCase());
					return count ? colorScale(count) : '#e5e7eb';
				})
				.attr('stroke', '#fff')
				.attr('stroke-width', 0.5)
				.style('cursor', 'pointer')
				.style('opacity', 0)
				.on('mouseover', function (event, d) {
					// Normalize country name (replace Israel with Palestine)
					const countryName = normalizeCountryName(d.properties.name);
					const count = countryDataMap.get(countryName?.toLowerCase());

					if (count) {
						d3.select(this)
							.transition()
							.duration(200)
							.attr('stroke-width', 2)
							.attr('stroke', '#4f46e5');

						tooltip.transition().duration(200).style('opacity', 1);

						tooltip
							.html(`<strong>${countryName}</strong><br/>Visitors: <strong>${count}</strong>`)
							.style('left', event.pageX + 10 + 'px')
							.style('top', event.pageY - 10 + 'px');
					}
				})
				.on('mouseout', function () {
					d3.select(this)
						.transition()
						.duration(200)
						.attr('stroke-width', 0.5)
						.attr('stroke', '#fff');

					tooltip.transition().duration(200).style('opacity', 0);
				})
				.on('click', function (event, d) {
					// Normalize country name (replace Israel with Palestine)
					const countryName = normalizeCountryName(d.properties.name);
					const count = countryDataMap.get(countryName?.toLowerCase());

					if (count && onclick) {
						onclick({ name: countryName, count });
					}
				})
				.transition()
				.duration(800)
				.delay((d, i) => i * 2)
				.style('opacity', 1);

			// Add legend
			const legendWidth = 200;
			const legendHeight = 10;
			const legendX = width - legendWidth - 20;
			const legendY = height - 40;

			const legendScale = d3.scaleLinear().domain([0, maxCount]).range([0, legendWidth]);

			const legendAxis = d3.axisBottom(legendScale).ticks(5).tickFormat(d3.format('.0f'));

			// Create gradient for legend
			const defs = svg.append('defs');
			const linearGradient = defs
				.append('linearGradient')
				.attr('id', 'legend-gradient')
				.attr('x1', '0%')
				.attr('x2', '100%');

			linearGradient
				.selectAll('stop')
				.data(d3.range(0, 1.1, 0.1))
				.join('stop')
				.attr('offset', (d) => d * 100 + '%')
				.attr('stop-color', (d) => colorScale(d * maxCount));

			// Draw legend
			const legend = svg.append('g').attr('transform', `translate(${legendX},${legendY})`);

			legend
				.append('rect')
				.attr('width', legendWidth)
				.attr('height', legendHeight)
				.style('fill', 'url(#legend-gradient)')
				.attr('stroke', '#ccc')
				.attr('stroke-width', 1);

			legend
				.append('g')
				.attr('transform', `translate(0,${legendHeight})`)
				.call(legendAxis)
				.selectAll('text')
				.style('font-size', '10px')
				.style('fill', '#6b7280');

			legend
				.append('text')
				.attr('x', legendWidth / 2)
				.attr('y', -5)
				.attr('text-anchor', 'middle')
				.style('font-size', '12px')
				.style('fill', '#374151')
				.text('Visitors');

			// Cleanup tooltip on destroy
			return () => {
				tooltip.remove();
			};
		} catch (error) {
			console.error('Error loading map data:', error);

			// Show error message
			svg
				.append('text')
				.attr('x', width / 2)
				.attr('y', height / 2)
				.attr('text-anchor', 'middle')
				.attr('fill', '#6b7280')
				.text('Error loading map data');
		}
	}

	onMount(() => {
		const cleanup = drawMap();

		// Handle window resize
		const handleResize = () => {
			drawMap();
		};

		window.addEventListener('resize', handleResize);

		return () => {
			window.removeEventListener('resize', handleResize);
			if (cleanup) cleanup();
		};
	});

	$effect(() => {
		// Redraw when data changes
		if (data) {
			drawMap();
		}
	});
</script>

<div bind:this={mapContainer} class="h-[400px] w-full">
	{#if data.length === 0}
		<div class="text-muted-foreground flex h-full items-center justify-center">
			<p>No geographical data available for the selected period</p>
		</div>
	{:else}
		<svg bind:this={svgElement} class="w-full"></svg>
	{/if}
</div>
