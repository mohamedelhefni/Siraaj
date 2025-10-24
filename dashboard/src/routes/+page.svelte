<script>
	import { onMount } from 'svelte';
	import { format, subDays } from 'date-fns';
	import { fetchStats } from '$lib/api';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import StatsCard from '$lib/components/StatsCard.svelte';
	import TimelineChart from '$lib/components/TimelineChart.svelte';
	import TopItemsList from '$lib/components/TopItemsList.svelte';
	import DateRangePicker from '$lib/components/DateRangePicker.svelte';

	let stats = $state({
		total_events: 0,
		unique_users: 0,
		top_events: [],
		timeline: [],
		top_pages: [],
		browsers: [],
		top_countries: [],
		top_sources: []
	});

	let loading = $state(true);
	let error = $state(null);

	// Default to last 7 days
	let startDate = $state(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	let endDate = $state(format(new Date(), 'yyyy-MM-dd'));

	async function loadStats() {
		loading = true;
		error = null;
		try {
			stats = await fetchStats(startDate, endDate);
		} catch (err) {
			error = err.message;
			console.error('Failed to load stats:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadStats();
	});

	function handleDateChange(event) {
		startDate = event.detail.startDate;
		endDate = event.detail.endDate;
		loadStats();
	}
</script>

<div class="container mx-auto space-y-6 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-4xl font-bold tracking-tight">Analytics Dashboard</h1>
			<p class="text-muted-foreground mt-2">Real-time insights and analytics powered by DuckDB</p>
		</div>
		<DateRangePicker {startDate} {endDate} on:change={handleDateChange} />
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="text-center">
				<div
					class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent motion-reduce:animate-[spin_1.5s_linear_infinite]"
				></div>
				<p class="text-muted-foreground mt-4">Loading analytics...</p>
			</div>
		</div>
	{:else if error}
		<Card class="border-destructive">
			<CardHeader>
				<CardTitle class="text-destructive">Error Loading Data</CardTitle>
				<CardDescription>{error}</CardDescription>
			</CardHeader>
		</Card>
	{:else}
		<!-- Overview Stats -->
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			<StatsCard
				title="Total Events"
				value={stats.total_events?.toLocaleString() || '0'}
				icon="activity"
				description="Total tracked events"
			/>
			<StatsCard
				title="Unique Users"
				value={stats.unique_users?.toLocaleString() || '0'}
				icon="users"
				description="Unique visitors"
			/>
			<StatsCard
				title="Avg Events/User"
				value={stats.unique_users > 0 ? (stats.total_events / stats.unique_users).toFixed(1) : '0'}
				icon="trending-up"
				description="Average engagement"
			/>
			<StatsCard
				title="Countries"
				value={stats.top_countries?.length || '0'}
				icon="globe"
				description="Geographic reach"
			/>
		</div>

		<!-- Timeline Chart -->
		<Card>
			<CardHeader>
				<CardTitle>Events Over Time</CardTitle>
				<CardDescription>Daily event tracking</CardDescription>
			</CardHeader>
			<CardContent>
				<TimelineChart data={stats.timeline || []} />
			</CardContent>
		</Card>

		<!-- Top Events and Pages -->
		<div class="grid gap-4 md:grid-cols-2">
			<Card>
				<CardHeader>
					<CardTitle>Top Events</CardTitle>
					<CardDescription>Most tracked events</CardDescription>
				</CardHeader>
				<CardContent>
					<TopItemsList items={stats.top_events || []} labelKey="name" valueKey="count" />
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Top Pages</CardTitle>
					<CardDescription>Most visited pages</CardDescription>
				</CardHeader>
				<CardContent>
					<TopItemsList items={stats.top_pages || []} labelKey="url" valueKey="count" />
				</CardContent>
			</Card>
		</div>

		<!-- Browsers, Countries, and Sources -->
		<div class="grid gap-4 md:grid-cols-3">
			<Card>
				<CardHeader>
					<CardTitle>Browsers</CardTitle>
					<CardDescription>Browser distribution</CardDescription>
				</CardHeader>
				<CardContent>
					<TopItemsList
						items={stats.browsers || []}
						labelKey="name"
						valueKey="count"
						maxItems={5}
					/>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Top Countries</CardTitle>
					<CardDescription>Geographic distribution</CardDescription>
				</CardHeader>
				<CardContent>
					<TopItemsList
						items={stats.top_countries || []}
						labelKey="name"
						valueKey="count"
						maxItems={5}
					/>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Top Sources</CardTitle>
					<CardDescription>Traffic sources</CardDescription>
				</CardHeader>
				<CardContent>
					<TopItemsList
						items={stats.top_sources || []}
						labelKey="name"
						valueKey="count"
						maxItems={5}
					/>
				</CardContent>
			</Card>
		</div>
	{/if}
</div>
