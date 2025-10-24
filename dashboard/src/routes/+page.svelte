<script>
	import { onMount, onDestroy } from 'svelte';
	import { format, subDays, startOfMonth, startOfYear, subMonths } from 'date-fns';
	import { fetchStats, fetchOnlineUsers, fetchProjects, fetchTopProperties } from '$lib/api';
	import { RefreshCw, X, TrendingUp, TrendingDown, Minus } from 'lucide-svelte';
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
	import PropertiesPanel from '$lib/components/PropertiesPanel.svelte';
	import CountriesPanel from '$lib/components/CountriesPanel.svelte';
	import DateRangePicker from '$lib/components/DateRangePicker.svelte';

	let stats = $state({
		total_events: 0,
		unique_users: 0,
		total_visits: 0,
		page_views: 0,
		bounce_rate: 0,
		events_change: 0,
		users_change: 0,
		visits_change: 0,
		page_views_change: 0,
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
	let topProperties = $state([]);

	let loading = $state(true);
	let statsLoading = $state(false);
	let propertiesLoading = $state(false);
	let onlineUsersLoading = $state(false);
	let error = $state(null);
	let autoRefresh = $state(true);
	let refreshInterval = $state(null);
	let refreshIntervalTime = $state(30000); // 30 seconds default
	let lastRefresh = $state(new Date());

	// Date range presets
	let dateRangePreset = $state('last_7_days');
	const dateRangePresets = [
		{ value: 'today', label: 'Today' },
		{ value: 'yesterday', label: 'Yesterday' },
		{ value: 'last_7_days', label: 'Last 7 days' },
		{ value: 'last_30_days', label: 'Last 30 days' },
		{ value: 'this_month', label: 'This month' },
		{ value: 'last_month', label: 'Last month' },
		{ value: 'last_3_months', label: 'Last 3 months' },
		{ value: 'last_6_months', label: 'Last 6 months' },
		{ value: 'this_year', label: 'This year' },
		{ value: 'custom', label: 'Custom range' }
	];

	// Filter states
	let activeFilters = $state({
		source: null,
		country: null,
		browser: null,
		event: null,
		project: null,
		metric: null, // For filtering by clicked metric card
		propertyKey: null,
		propertyValue: null
	});

	// Default to last 7 days
	let startDate = $state(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	let endDate = $state(format(new Date(), 'yyyy-MM-dd'));
	let showCustomDateInputs = $state(false);

	// Apply date range preset
	function applyDateRangePreset(preset) {
		const now = new Date();
		const today = format(now, 'yyyy-MM-dd');

		switch (preset) {
			case 'today':
				startDate = today;
				endDate = today;
				break;
			case 'yesterday':
				const yesterday = subDays(now, 1);
				startDate = format(yesterday, 'yyyy-MM-dd');
				endDate = format(yesterday, 'yyyy-MM-dd');
				break;
			case 'last_7_days':
				startDate = format(subDays(now, 7), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'last_30_days':
				startDate = format(subDays(now, 30), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'this_month':
				startDate = format(startOfMonth(now), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'last_month':
				const lastMonthStart = startOfMonth(subMonths(now, 1));
				const lastMonthEnd = subDays(startOfMonth(now), 1);
				startDate = format(lastMonthStart, 'yyyy-MM-dd');
				endDate = format(lastMonthEnd, 'yyyy-MM-dd');
				break;
			case 'last_3_months':
				startDate = format(subMonths(now, 3), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'last_6_months':
				startDate = format(subMonths(now, 6), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'this_year':
				startDate = format(startOfYear(now), 'yyyy-MM-dd');
				endDate = today;
				break;
			case 'custom':
				showCustomDateInputs = true;
				return;
		}
		showCustomDateInputs = false;
		loadStats();
	}

	// URL param helpers
	function updateURLParams() {
		if (typeof window === 'undefined') return;

		const params = new URLSearchParams();
		params.set('range', dateRangePreset);
		if (dateRangePreset === 'custom') {
			params.set('start', startDate);
			params.set('end', endDate);
		}
		if (activeFilters.project) params.set('project', activeFilters.project);
		if (activeFilters.source) params.set('source', activeFilters.source);
		if (activeFilters.country) params.set('country', activeFilters.country);
		if (activeFilters.browser) params.set('browser', activeFilters.browser);
		if (activeFilters.event) params.set('event', activeFilters.event);
		if (activeFilters.metric) params.set('metric', activeFilters.metric);
		if (activeFilters.propertyKey) params.set('propKey', activeFilters.propertyKey);
		if (activeFilters.propertyValue) params.set('propValue', activeFilters.propertyValue);
		if (refreshIntervalTime !== 30000) params.set('interval', refreshIntervalTime.toString());

		const newURL = `${window.location.pathname}?${params.toString()}`;
		window.history.replaceState({}, '', newURL);
	}
	function loadFromURLParams() {
		if (typeof window === 'undefined') return;

		const params = new URLSearchParams(window.location.search);

		if (params.has('range')) {
			dateRangePreset = params.get('range');
			if (dateRangePreset === 'custom') {
				if (params.has('start')) startDate = params.get('start');
				if (params.has('end')) endDate = params.get('end');
				showCustomDateInputs = true;
			} else {
				applyDateRangePreset(dateRangePreset);
			}
		}
		if (params.has('project')) activeFilters.project = params.get('project');
		if (params.has('source')) activeFilters.source = params.get('source');
		if (params.has('country')) activeFilters.country = params.get('country');
		if (params.has('browser')) activeFilters.browser = params.get('browser');
		if (params.has('event')) activeFilters.event = params.get('event');
		if (params.has('metric')) activeFilters.metric = params.get('metric');
		if (params.has('propKey')) activeFilters.propertyKey = params.get('propKey');
		if (params.has('propValue')) activeFilters.propertyValue = params.get('propValue');
		if (params.has('interval')) {
			refreshIntervalTime = parseInt(params.get('interval'));
		}
	}

	async function loadStats() {
		if (loading) {
			// First time loading - show full page loader
			loading = true;
		} else {
			// Subsequent loads - show component-level loaders
			statsLoading = true;
			propertiesLoading = true;
			onlineUsersLoading = true;
		}

		error = null;

		try {
			// Load stats
			const statsPromise = fetchStats(startDate, endDate, 50, activeFilters)
				.then((data) => {
					stats = data;
					statsLoading = false;
					return data;
				})
				.catch((err) => {
					statsLoading = false;
					throw err;
				});

			// Load online users
			const onlinePromise = fetchOnlineUsers(5)
				.then((data) => {
					onlineData = data;
					onlineUsersLoading = false;
					return data;
				})
				.catch((err) => {
					onlineUsersLoading = false;
					throw err;
				});

			// Load properties
			const propertiesPromise = fetchTopProperties(startDate, endDate, 5, activeFilters)
				.then((data) => {
					topProperties = data || [];
					propertiesLoading = false;
					return data;
				})
				.catch((err) => {
					propertiesLoading = false;
					throw err;
				});

			await Promise.all([statsPromise, onlinePromise, propertiesPromise]);

			lastRefresh = new Date();
			updateURLParams();
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
		updateURLParams();
	}

	function addFilter(type, value) {
		activeFilters[type] = value;
		updateURLParams();
		loadStats();
	}

	function removeFilter(type) {
		activeFilters[type] = null;
		updateURLParams();
		loadStats();
	}

	function clearAllFilters() {
		activeFilters = {
			source: null,
			country: null,
			browser: null,
			event: null,
			project: null,
			metric: null,
			propertyKey: null,
			propertyValue: null
		};
		updateURLParams();
		loadStats();
	}

	onMount(() => {
		loadFromURLParams();
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
		updateURLParams();
		loadStats();
	}

	// Helper to format trend
	function getTrendIcon(change) {
		if (change > 0) return TrendingUp;
		if (change < 0) return TrendingDown;
		return Minus;
	}

	function getTrendColor(change) {
		if (change > 0) return 'text-green-600';
		if (change < 0) return 'text-red-600';
		return 'text-gray-600';
	}

	// Handle metric card clicks for filtering timeline
	function handleMetricClick(metricType) {
		if (activeFilters.metric === metricType) {
			// Toggle off if already selected
			activeFilters.metric = null;
		} else {
			activeFilters.metric = metricType;
		}
		updateURLParams();
		loadStats(); // Reload data with metric filter
	}

	// Check if a metric is selected
	function isMetricSelected(metricType) {
		return activeFilters.metric === metricType;
	}

	// Handle property filter
	function addPropertyFilter(prop) {
		activeFilters.propertyKey = prop.key;
		activeFilters.propertyValue = prop.value;
		updateURLParams();
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
	</div>

	<!-- Controls Row -->
	<div class="flex flex-wrap items-center gap-3">
		<!-- Date Range Selector -->
		<div class="flex items-center gap-2">
			<label class="text-sm font-medium">Period:</label>
			<select
				class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
				bind:value={dateRangePreset}
				onchange={(e) => applyDateRangePreset(e.target.value)}
			>
				{#each dateRangePresets as preset}
					<option value={preset.value}>{preset.label}</option>
				{/each}
			</select>
		</div>

		<!-- Custom Date Inputs (shown when custom is selected) -->
		{#if showCustomDateInputs}
			<div class="flex items-center gap-2 rounded-lg border p-2">
				<input
					type="date"
					bind:value={startDate}
					class="border-none text-sm focus:outline-none"
					onchange={() => loadStats()}
				/>
				<span class="text-muted-foreground">to</span>
				<input
					type="date"
					bind:value={endDate}
					class="border-none text-sm focus:outline-none"
					onchange={() => loadStats()}
				/>
			</div>
		{/if}

		<!-- Project Selector -->
		{#if projects.length > 0}
			<div class="flex items-center gap-2">
				<label class="text-sm font-medium">Project:</label>
				<select
					class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
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
			</div>
		{/if}

		<!-- Auto-refresh controls -->
		<div class="ml-auto flex items-center gap-2">
			<label class="text-sm font-medium">Auto-refresh:</label>
			<select
				class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
				value={refreshIntervalTime}
				onchange={handleRefreshIntervalChange}
			>
				<option value="0">Off</option>
				<option value="10000">10s</option>
				<option value="30000">30s</option>
				<option value="60000">1min</option>
				<option value="300000">5min</option>
			</select>

			<Button variant="outline" size="sm" onclick={() => loadStats()} class="gap-2">
				<RefreshCw class="h-4 w-4" />
				Refresh
			</Button>

			{#if !loading}
				<span class="text-muted-foreground whitespace-nowrap text-xs">
					Updated {format(lastRefresh, 'HH:mm:ss')}
				</span>
			{/if}
		</div>
	</div>

	<!-- Active Filters -->
	{#if Object.values(activeFilters).some((v) => v !== null)}
		<div class="flex flex-wrap items-center gap-2">
			<span class="text-muted-foreground text-sm">Active Filters:</span>
			{#if activeFilters.project}
				<Badge variant="secondary" class="gap-1">
					Project: {activeFilters.project}
					<button onclick={() => removeFilter('project')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.source}
				<Badge variant="secondary" class="gap-1">
					Source: {activeFilters.source}
					<button onclick={() => removeFilter('source')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.country}
				<Badge variant="secondary" class="gap-1">
					Country: {activeFilters.country}
					<button onclick={() => removeFilter('country')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.browser}
				<Badge variant="secondary" class="gap-1">
					Browser: {activeFilters.browser}
					<button onclick={() => removeFilter('browser')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.event}
				<Badge variant="secondary" class="gap-1">
					Event: {activeFilters.event}
					<button onclick={() => removeFilter('event')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.metric}
				<Badge variant="secondary" class="gap-1">
					Metric: {activeFilters.metric}
					<button onclick={() => removeFilter('metric')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.propertyKey && activeFilters.propertyValue}
				<Badge variant="secondary" class="gap-1">
					Property: {activeFilters.propertyKey}={activeFilters.propertyValue}
					<button
						onclick={() => {
							removeFilter('propertyKey');
							removeFilter('propertyValue');
						}}
						class="hover:text-destructive ml-1"
					>
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			<Button variant="ghost" size="sm" onclick={clearAllFilters}>Clear All</Button>
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
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-6">
			<Card
				class="cursor-pointer transition-all hover:shadow-md {isMetricSelected('events')
					? 'ring-primary ring-2'
					: ''}"
				onclick={() => handleMetricClick('events')}
			>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Total Events</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stats.total_events?.toLocaleString() || '0'}</div>
					{#if stats.events_change !== undefined && stats.events_change !== 0}
						{@const TrendIcon = getTrendIcon(stats.events_change)}
						<p class="text-xs {getTrendColor(stats.events_change)} mt-1 flex items-center gap-1">
							<TrendIcon class="h-3 w-3" />
							{Math.abs(stats.events_change).toFixed(1)}% vs previous period
						</p>
					{:else}
						<p class="text-muted-foreground mt-1 text-xs">All tracked events</p>
					{/if}
				</CardContent>
			</Card>

			<Card
				class="cursor-pointer transition-all hover:shadow-md {isMetricSelected('users')
					? 'ring-primary ring-2'
					: ''}"
				onclick={() => handleMetricClick('users')}
			>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Unique Visitors</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stats.unique_users?.toLocaleString() || '0'}</div>
					{#if stats.users_change !== undefined && stats.users_change !== 0}
						{@const TrendIcon = getTrendIcon(stats.users_change)}
						<p class="text-xs {getTrendColor(stats.users_change)} mt-1 flex items-center gap-1">
							<TrendIcon class="h-3 w-3" />
							{Math.abs(stats.users_change).toFixed(1)}% vs previous period
						</p>
					{:else}
						<p class="text-muted-foreground mt-1 text-xs">Unique users</p>
					{/if}
				</CardContent>
			</Card>

			<Card
				class="cursor-pointer transition-all hover:shadow-md {isMetricSelected('visits')
					? 'ring-primary ring-2'
					: ''}"
				onclick={() => handleMetricClick('visits')}
			>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Total Visits</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stats.total_visits?.toLocaleString() || '0'}</div>
					{#if stats.visits_change !== undefined && stats.visits_change !== 0}
						{@const TrendIcon = getTrendIcon(stats.visits_change)}
						<p class="text-xs {getTrendColor(stats.visits_change)} mt-1 flex items-center gap-1">
							<TrendIcon class="h-3 w-3" />
							{Math.abs(stats.visits_change).toFixed(1)}% vs previous period
						</p>
					{:else}
						<p class="text-muted-foreground mt-1 text-xs">Unique sessions</p>
					{/if}
				</CardContent>
			</Card>

			<Card
				class="cursor-pointer transition-all hover:shadow-md {isMetricSelected('page_views')
					? 'ring-primary ring-2'
					: ''}"
				onclick={() => handleMetricClick('page_views')}
			>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Page Views</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stats.page_views?.toLocaleString() || '0'}</div>
					{#if stats.page_views_change !== undefined && stats.page_views_change !== 0}
						{@const TrendIcon = getTrendIcon(stats.page_views_change)}
						<p
							class="text-xs {getTrendColor(stats.page_views_change)} mt-1 flex items-center gap-1"
						>
							<TrendIcon class="h-3 w-3" />
							{Math.abs(stats.page_views_change).toFixed(1)}% vs previous period
						</p>
					{:else}
						<p class="text-muted-foreground mt-1 text-xs">Total page views</p>
					{/if}
				</CardContent>
			</Card>

			<Card
				class=" transition-all hover:shadow-md {isMetricSelected('bounce_rate')
					? 'ring-primary ring-2'
					: ''}"
			>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Bounce Rate</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stats.bounce_rate?.toFixed(1) || '0'}%</div>
					<p class="text-muted-foreground mt-1 text-xs">Single page sessions</p>
				</CardContent>
			</Card>

			<Card>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">Online Now</CardTitle>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{onlineData.online_users?.toLocaleString() || '0'}</div>
					<p class="text-muted-foreground mt-1 text-xs">Active in last 5 min</p>
				</CardContent>
			</Card>
		</div>

		<!-- Timeline Chart -->
		<Card>
			<CardHeader>
				<CardTitle>
					{#if activeFilters.metric === 'users'}
						Unique Visitors Over Time
					{:else if activeFilters.metric === 'visits'}
						Total Visits Over Time
					{:else if activeFilters.metric === 'page_views'}
						Page Views Over Time
					{:else}
						Events Over Time
					{/if}
				</CardTitle>
				<CardDescription>
					{#if stats.timeline_format === 'hour'}
						Hourly tracking
					{:else if stats.timeline_format === 'month'}
						Monthly tracking
					{:else}
						Daily tracking
					{/if}
					{#if activeFilters.metric}
						· Click the metric card again to show all events
					{/if}
				</CardDescription>
			</CardHeader>
			<CardContent>
				{#if statsLoading}
					<div class="text-muted-foreground flex h-[300px] items-center justify-center">
						<div class="flex flex-col items-center gap-2">
							<div
								class="border-primary h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
							></div>
							<p class="text-sm">Loading timeline...</p>
						</div>
					</div>
				{:else}
					<TimelineChart
						data={stats.timeline || []}
						format={stats.timeline_format || 'day'}
						metric={activeFilters.metric || 'events'}
					/>
				{/if}
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
					{#if statsLoading}
						<div class="text-muted-foreground flex min-h-[150px] items-center justify-center">
							<div class="flex flex-col items-center gap-2">
								<div
									class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
								></div>
								<p class="text-xs">Loading...</p>
							</div>
						</div>
					{:else}
						<TopItemsList
							items={stats.top_events || []}
							labelKey="name"
							valueKey="count"
							maxItems={5}
							showMoreTitle="All Events ({(stats.top_events || []).length} total)"
							onclick={(item) => addFilter('event', item.name)}
						/>
					{/if}
				</CardContent>
			</Card>
			<Card>
				<CardHeader>
					<CardTitle>Top Countries</CardTitle>
					<CardDescription>
						Geographic distribution
						{#if stats.top_countries && stats.top_countries.length > 0}
							· {stats.top_countries.length} countries
						{/if}
					</CardDescription>
				</CardHeader>
				<CardContent>
					<CountriesPanel
						countries={stats.top_countries || []}
						onclick={(item) => addFilter('country', item.name)}
						loading={statsLoading}
					/>
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
					{#if statsLoading}
						<div class="text-muted-foreground flex min-h-[150px] items-center justify-center">
							<div class="flex flex-col items-center gap-2">
								<div
									class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
								></div>
								<p class="text-xs">Loading...</p>
							</div>
						</div>
					{:else}
						<TopItemsList
							items={stats.browsers || []}
							labelKey="name"
							valueKey="count"
							maxItems={5}
							type="browser"
							showMoreTitle="All Browsers ({(stats.browsers || []).length} total)"
							onclick={(item) => addFilter('browser', item.name)}
						/>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Top Pages</CardTitle>
					<CardDescription>Most visited pages</CardDescription>
				</CardHeader>
				<CardContent>
					{#if statsLoading}
						<div class="text-muted-foreground flex min-h-[150px] items-center justify-center">
							<div class="flex flex-col items-center gap-2">
								<div
									class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
								></div>
								<p class="text-xs">Loading...</p>
							</div>
						</div>
					{:else}
						<TopItemsList
							items={stats.top_pages || []}
							labelKey="url"
							maxItems={5}
							valueKey="count"
							showMoreTitle="All Pages ({(stats.top_pages || []).length} total)"
						/>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Top Sources</CardTitle>
					<CardDescription>Traffic sources</CardDescription>
				</CardHeader>
				<CardContent>
					{#if statsLoading}
						<div class="text-muted-foreground flex min-h-[150px] items-center justify-center">
							<div class="flex flex-col items-center gap-2">
								<div
									class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
								></div>
								<p class="text-xs">Loading...</p>
							</div>
						</div>
					{:else}
						<TopItemsList
							items={stats.top_sources || []}
							labelKey="name"
							valueKey="count"
							maxItems={5}
							type="source"
							showMoreTitle="All Sources ({(stats.top_sources || []).length} total)"
							onclick={(item) => addFilter('source', item.name)}
						/>
					{/if}
				</CardContent>
			</Card>
		</div>

		<!-- Custom Properties Panel -->
		<Card>
			<CardHeader>
				<CardTitle>Custom Event Properties</CardTitle>
				<CardDescription>
					Top property key-value pairs from tracked events
					{#if topProperties.length > 0}
						· {topProperties.length} unique properties
					{/if}
				</CardDescription>
			</CardHeader>
			<CardContent>
				{#if propertiesLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
						<div class="flex flex-col items-center gap-2">
							<div
								class="border-primary h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
							></div>
							<p class="text-sm">Loading properties...</p>
						</div>
					</div>
				{:else}
					<PropertiesPanel properties={topProperties} onPropertyClick={addPropertyFilter} />
				{/if}
			</CardContent>
		</Card>
	{/if}
</div>
