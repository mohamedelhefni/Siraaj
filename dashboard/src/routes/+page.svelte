<script>
	import { onMount, onDestroy } from 'svelte';
	import { format, subDays } from 'date-fns';
	import { fetchStats, fetchOnlineUsers, fetchProjects } from '$lib/api';
	import { RefreshCw, X } from 'lucide-svelte';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
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

	let onlineData = $state({
		online_users: 0,
		active_sessions: 0
	});

	let projects = $state([]);

	let loading = $state(true);
	let error = $state(null);
	let autoRefresh = $state(true);
	let refreshInterval = $state(null);
	let refreshIntervalTime = $state(30000); // 30 seconds default
	let lastRefresh = $state(new Date());

	// Filter states
	let activeFilters = $state({
		source: null,
		country: null,
		browser: null,
		event: null,
		project: null
	});

	// Default to last 7 days
	let startDate = $state(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	let endDate = $state(format(new Date(), 'yyyy-MM-dd'));

	async function loadStats() {
		loading = true;
		error = null;
		try {
			const [statsData, onlineUsersData] = await Promise.all([
				fetchStats(startDate, endDate, 50, activeFilters),
				fetchOnlineUsers(5)
			]);
			stats = statsData;
			onlineData = onlineUsersData;
			lastRefresh = new Date();
		} catch (err) {
			error = err.message;
			console.error('Failed to load stats:', err);
		} finally {
			loading = false;
		}
	}

	function setupAutoRefresh() {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}

		if (autoRefresh && refreshIntervalTime > 0) {
			refreshInterval = setInterval(() => {
				loadStats();
			}, refreshIntervalTime);
		}
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		setupAutoRefresh();
	}

	function handleRefreshIntervalChange(event) {
		refreshIntervalTime = parseInt(event.target.value);
		setupAutoRefresh();
	}

	function addFilter(type, value) {
		activeFilters[type] = value;
		loadStats();
	}

	function removeFilter(type) {
		activeFilters[type] = null;
		loadStats();
	}

	function clearAllFilters() {
		activeFilters = {
			source: null,
			country: null,
			browser: null,
			event: null,
			project: null
		};
		loadStats();
	}

	onMount(() => {
		loadProjects();
		loadStats();
		setupAutoRefresh();
	});

	onDestroy(() => {
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	});

	async function loadProjects() {
		try {
			projects = await fetchProjects();
		} catch (err) {
			console.error('Failed to load projects:', err);
		}
	}

	function handleDateChange(event) {
		startDate = event.detail.startDate;
		endDate = event.detail.endDate;
		loadStats();
	}
</script>

<div class="container mx-auto space-y-6 p-6">
	<!-- Header -->
	<div class="flex flex-wrap items-center justify-between gap-4">
		<div>
			<h1 class="text-4xl font-bold tracking-tight">Siraaj Dashboard</h1>
			<p class="text-muted-foreground mt-2">Real-time insights and analytics powered by DuckDB</p>
		</div>
		<div class="flex flex-wrap items-center gap-2">
			<DateRangePicker {startDate} {endDate} on:change={handleDateChange} />
			
			<!-- Project Selector -->
			{#if projects.length > 0}
				<select
					class="flex h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
					value={activeFilters.project || ''}
					onchange={(e) => {
						if (e.target.value) {
							addFilter('project', e.target.value);
						} else {
							removeFilter('project');
						}
					}}
				>
					<option value="">All Projects</option>
					{#each projects as project}
						<option value={project}>{project}</option>
					{/each}
				</select>
			{/if}

			<!-- Auto-refresh controls -->
			<div class="flex items-center gap-2 rounded-lg border p-2">
				<Button
					variant={autoRefresh ? 'default' : 'outline'}
					size="sm"
					onclick={toggleAutoRefresh}
					class="gap-2"
				>
					<RefreshCw class={autoRefresh ? 'h-4 w-4 animate-spin' : 'h-4 w-4'} />
					{autoRefresh ? 'Auto' : 'Manual'}
				</Button>
				<select
					class="h-9 rounded-md border border-input bg-background px-2 py-1 text-xs"
					value={refreshIntervalTime}
					onchange={handleRefreshIntervalChange}
				>
					<option value="10000">10s</option>
					<option value="30000">30s</option>
					<option value="60000">1min</option>
					<option value="300000">5min</option>
					<option value="0">Off</option>
				</select>
				{#if !loading}
					<span class="text-muted-foreground text-xs whitespace-nowrap">
						{format(lastRefresh, 'HH:mm:ss')}
					</span>
				{/if}
			</div>
		</div>
	</div>

	<!-- Active Filters -->
	{#if Object.values(activeFilters).some(v => v !== null)}
		<div class="flex flex-wrap items-center gap-2">
			<span class="text-sm text-muted-foreground">Active Filters:</span>
			{#if activeFilters.project}
				<Badge variant="secondary" class="gap-1">
					Project: {activeFilters.project}
					<button onclick={() => removeFilter('project')} class="ml-1 hover:text-destructive">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.source}
				<Badge variant="secondary" class="gap-1">
					Source: {activeFilters.source}
					<button onclick={() => removeFilter('source')} class="ml-1 hover:text-destructive">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.country}
				<Badge variant="secondary" class="gap-1">
					Country: {activeFilters.country}
					<button onclick={() => removeFilter('country')} class="ml-1 hover:text-destructive">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.browser}
				<Badge variant="secondary" class="gap-1">
					Browser: {activeFilters.browser}
					<button onclick={() => removeFilter('browser')} class="ml-1 hover:text-destructive">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.event}
				<Badge variant="secondary" class="gap-1">
					Event: {activeFilters.event}
					<button onclick={() => removeFilter('event')} class="ml-1 hover:text-destructive">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			<Button variant="ghost" size="sm" onclick={clearAllFilters}>
				Clear All
			</Button>
		</div>
	{/if}

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
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
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
				title="Online Now"
				value={onlineData.online_users?.toLocaleString() || '0'}
				icon="eye"
				description="Active in last 5 min"
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
					<TopItemsList 
						items={stats.top_events || []} 
						labelKey="name" 
						valueKey="count"
						onclick={(item) => addFilter('event', item.name)}
					/>
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
						type="browser"
						onclick={(item) => addFilter('browser', item.name)}
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
						type="country"
						onclick={(item) => addFilter('country', item.name)}
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
						type="source"
						onclick={(item) => addFilter('source', item.name)}
					/>
				</CardContent>
			</Card>
		</div>
	{/if}
</div>
