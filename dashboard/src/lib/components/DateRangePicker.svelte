<script>
	import { createEventDispatcher } from 'svelte';
	import { Calendar } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';

	let { startDate, endDate } = $props();
	const dispatch = createEventDispatcher();

	let localStartDate = $state(startDate);
	let localEndDate = $state(endDate);

	function handleApply() {
		dispatch('change', {
			startDate: localStartDate,
			endDate: localEndDate
		});
	}

	function setLast7Days() {
		const end = new Date();
		const start = new Date();
		start.setDate(start.getDate() - 7);
		localStartDate = start.toISOString().split('T')[0];
		localEndDate = end.toISOString().split('T')[0];
		handleApply();
	}

	function setLast30Days() {
		const end = new Date();
		const start = new Date();
		start.setDate(start.getDate() - 30);
		localStartDate = start.toISOString().split('T')[0];
		localEndDate = end.toISOString().split('T')[0];
		handleApply();
	}

	function setThisMonth() {
		const now = new Date();
		const start = new Date(now.getFullYear(), now.getMonth(), 1);
		const end = new Date();
		localStartDate = start.toISOString().split('T')[0];
		localEndDate = end.toISOString().split('T')[0];
		handleApply();
	}
</script>

<div class="flex items-center gap-2 flex-wrap">
	<div class="flex items-center gap-2 border rounded-lg p-2">
		<Calendar class="h-4 w-4 text-muted-foreground" />
		<input
			type="date"
			bind:value={localStartDate}
			class="text-sm border-none focus:outline-none"
		/>
		<span class="text-muted-foreground">to</span>
		<input type="date" bind:value={localEndDate} class="text-sm border-none focus:outline-none" />
		<Button size="sm" onclick={handleApply}>Apply</Button>
	</div>

	<div class="flex gap-2">
		<Button variant="outline" size="sm" onclick={setLast7Days}>Last 7 days</Button>
		<Button variant="outline" size="sm" onclick={setLast30Days}>Last 30 days</Button>
		<Button variant="outline" size="sm" onclick={setThisMonth}>This month</Button>
	</div>
</div>
