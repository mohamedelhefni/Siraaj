<script>
	import { onMount } from 'svelte';
	import { Chart, registerables } from 'chart.js';

	let { data = [] } = $props();
	let canvas = $state();
	let chart = $state();

	Chart.register(...registerables);

	onMount(() => {
		if (canvas) {
			chart = new Chart(canvas, {
				type: 'line',
				data: {
					labels: data.map((d) => d.date),
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
							intersect: false
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
			chart.data.labels = data.map((d) => d.date);
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
