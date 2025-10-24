<script>
	import { onMount } from 'svelte';
	import { Chart, registerables } from 'chart.js';

	let { data = [], format = 'day' } = $props();
	let canvas = $state();
	let chart = $state();

	Chart.register(...registerables);

	// Format date labels based on granularity
	function formatLabel(dateStr) {
		const date = new Date(dateStr);

		if (format === 'hour') {
			// For hourly: show time (e.g., "14:00", "15:00")
			return date.toLocaleTimeString('en-US', {
				hour: '2-digit',
				minute: '2-digit',
				hour12: false
			});
		} else if (format === 'month') {
			// For monthly: show month and year (e.g., "Jan 2025")
			return date.toLocaleDateString('en-US', {
				month: 'short',
				year: 'numeric'
			});
		} else {
			// For daily: show date (e.g., "Oct 24")
			return date.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric'
			});
		}
	}

	onMount(() => {
		if (canvas) {
			chart = new Chart(canvas, {
				type: 'line',
				data: {
					labels: data.map((d) => formatLabel(d.date)),
					datasets: [
						{
							label: 'Events',
							data: data.map((d) => d.count),
							borderColor: 'rgb(99, 102, 241)',
							backgroundColor: 'rgba(99, 102, 241, 0.1)',
							tension: 0.4,
							fill: true
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					plugins: {
						legend: {
							display: false
						},
						tooltip: {
							mode: 'index',
							intersect: false,
							callbacks: {
								title: function (context) {
									// Show full date in tooltip
									const index = context[0].dataIndex;
									const dateStr = data[index]?.date;
									if (!dateStr) return '';

									const date = new Date(dateStr);
									if (format === 'hour') {
										return date.toLocaleString('en-US', {
											month: 'short',
											day: 'numeric',
											hour: '2-digit',
											minute: '2-digit',
											hour12: false
										});
									} else if (format === 'month') {
										return date.toLocaleDateString('en-US', {
											month: 'long',
											year: 'numeric'
										});
									} else {
										return date.toLocaleDateString('en-US', {
											month: 'long',
											day: 'numeric',
											year: 'numeric'
										});
									}
								}
							}
						}
					},
					scales: {
						y: {
							beginAtZero: true,
							grid: {
								color: 'rgba(0, 0, 0, 0.05)'
							}
						},
						x: {
							grid: {
								display: false
							}
						}
					}
				}
			});
		}

		return () => {
			if (chart) {
				chart.destroy();
			}
		};
	});

	$effect(() => {
		if (chart && data) {
			chart.data.labels = data.map((d) => formatLabel(d.date));
			chart.data.datasets[0].data = data.map((d) => d.count);
			chart.update();
		}
	});
</script>

<div class="h-[300px]">
	{#if data.length === 0}
		<div class="text-muted-foreground flex h-full items-center justify-center">
			<p>No data available for the selected period</p>
		</div>
	{:else}
		<canvas bind:this={canvas}></canvas>
	{/if}
</div>
