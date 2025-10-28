<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { format, subDays, startOfMonth, startOfYear, subMonths } from 'date-fns';
	import {
		fetchOnlineUsers,
		fetchProjects,
		fetchTopStats,
		fetchTimeline,
		fetchTopPages,
		fetchEntryExitPages,
		fetchTopCountries,
		fetchTopSources,
		fetchTopEvents,
		fetchBrowsersDevicesOS
	} from '$lib/api';
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
	import BrowserPanel from '$lib/components/BrowserPanel.svelte';
	import DateRangePicker from '$lib/components/DateRangePicker.svelte';
	import MetricCard from '$lib/components/MetricCard.svelte';

	let stats: any = $state({
		total_events: 0,
		unique_users: 0,
		total_visits: 0,
		page_views: 0,
		bounce_rate: 0,
		avg_session_duration: 0,
		bot_events: 0,
		human_events: 0,
		bot_users: 0,
		human_users: 0,
		bot_percentage: 0,
		events_change: 0,
		users_change: 0,
		visits_change: 0,
		page_views_change: 0
	});

	let timeline: any = $state({
		timeline: [],
		timeline_format: 'day'
	});

	let topPages: any = $state({
		top_pages: []
	});

	let entryExitPages: any = $state({
		entry_pages: [],
		exit_pages: []
	});

	let topCountries: any[] = $state([]);
	let topSources: any[] = $state([]);
	let topEvents: any[] = $state([]);

	let browsersDevicesOS: any = $state({
		browsers: [],
		devices: [],
		os: []
	});

	let comparisonStats: any = $state({
		unique_users: 0,
		total_visits: 0,
		page_views: 0,
		total_events: 0,
		bounce_rate: 0,
		avg_session_duration: 0,
		bot_percentage: 0
	});

	let onlineData = $state({
		online_users: 0,
		active_sessions: 0
	});

	let projects: string[] = $state([]);
	let topProperties: any[] = $state([]);

	let loading = $state(true);
	let statsLoading = $state(false);
	let propertiesLoading = $state(false);
	let onlineUsersLoading = $state(false);
	let error = $state<string | null>(null);
	let autoRefresh = $state(true);
	let refreshInterval = $state<number | null>(null);
	let refreshIntervalTime = $state(30000); // 30 seconds default
	let lastRefresh = $state(new Date());

	// Pages tab state
	let pagesTab = $state<'all' | 'entry' | 'exit'>('all');

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
	let activeFilters = $state<{
		source: string | null;
		country: string | null;
		browser: string | null;
		device: string | null;
		os: string | null;
		event: string | null;
		project: string | null;
		metric: string | null;
		propertyKey: string | null;
		propertyValue: string | null;
		botFilter: string | null; // 'human', 'bot', or null for all
		page: string | null;
	}>({
		source: null,
		country: null,
		browser: null,
		device: null,
		os: null,
		event: null,
		project: null,
		metric: null, // For filtering by clicked metric card
		propertyKey: null,
		propertyValue: null,
		botFilter: null,
		page: null
	});

	// Default to last 7 days
	let startDate = $state(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	let endDate = $state(format(new Date(), 'yyyy-MM-dd'));
	let showCustomDateInputs = $state(false);

	// Computed period labels for display
	const currentPeriodLabel = $derived(() => {
		const start = new Date(startDate);
		const end = new Date(endDate);
		return `${format(start, 'd MMM')} - ${format(end, 'd MMM')}`;
	});

	const previousPeriodLabel = $derived(() => {
		const start = new Date(startDate);
		const end = new Date(endDate);
		const duration = end.getTime() - start.getTime();
		const prevEnd = new Date(start.getTime() - 24 * 60 * 60 * 1000);
		const prevStart = new Date(prevEnd.getTime() - duration);
		return `${format(prevStart, 'd MMM')} - ${format(prevEnd, 'd MMM')}`;
	});

	// Apply date range preset
	function applyDateRangePreset(preset: string) {
		applyDateRangePresetWithoutLoad(preset);
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
		if (activeFilters.device) params.set('device', activeFilters.device);
		if (activeFilters.os) params.set('os', activeFilters.os);
		if (activeFilters.event) params.set('event', activeFilters.event);
		if (activeFilters.metric) params.set('metric', activeFilters.metric);
		if (activeFilters.propertyKey) params.set('propKey', activeFilters.propertyKey);
		if (activeFilters.propertyValue) params.set('propValue', activeFilters.propertyValue);
		if (activeFilters.botFilter) params.set('botFilter', activeFilters.botFilter);
		if (activeFilters.page) params.set('page', activeFilters.page);
		if (refreshIntervalTime !== 30000) params.set('interval', refreshIntervalTime.toString());

		const newURL = `${window.location.pathname}?${params.toString()}`;
		window.history.replaceState({}, '', newURL);
	}
	function loadFromURLParams() {
		if (typeof window === 'undefined') return;

		const params = new URLSearchParams(window.location.search);

		if (params.has('range')) {
			const range = params.get('range');
			if (range) dateRangePreset = range;
			if (dateRangePreset === 'custom') {
				const start = params.get('start');
				const end = params.get('end');
				if (start) startDate = start;
				if (end) endDate = end;
				showCustomDateInputs = true;
			} else {
				// Don't call loadStats() here - will be called in onMount
				applyDateRangePresetWithoutLoad(dateRangePreset);
			}
		}
		const project = params.get('project');
		const source = params.get('source');
		const country = params.get('country');
		const browser = params.get('browser');
		const device = params.get('device');
		const os = params.get('os');
		const event = params.get('event');
		const metric = params.get('metric');
		const propKey = params.get('propKey');
		const propValue = params.get('propValue');
		const botFilter = params.get('botFilter');
		const page = params.get('page');
		const interval = params.get('interval');

		if (project) activeFilters.project = project;
		if (source) activeFilters.source = source;
		if (country) activeFilters.country = country;
		if (browser) activeFilters.browser = browser;
		if (device) activeFilters.device = device;
		if (os) activeFilters.os = os;
		if (event) activeFilters.event = event;
		if (metric) activeFilters.metric = metric;
		if (propKey) activeFilters.propertyKey = propKey;
		if (propValue) activeFilters.propertyValue = propValue;
		if (botFilter) activeFilters.botFilter = botFilter;
		if (page) activeFilters.page = page;
		if (interval) {
			refreshIntervalTime = parseInt(interval);
		}
	}

	// Apply date range preset without triggering loadStats (for initial load)
	function applyDateRangePresetWithoutLoad(preset: string) {
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
			// Load all stats endpoints in parallel for maximum performance
			const [
				topStatsData,
				timelineData,
				pagesData,
				entryExitData,
				countriesData,
				sourcesData,
				eventsData,
				devicesData,
				onlineData,
				comparisonData
			] = await Promise.all([
				// Main stats (counts, rates, trends)
				fetchTopStats(startDate, endDate, activeFilters).catch((err) => {
					console.error('Failed to load top stats:', err);
					return null;
				}),

				// Timeline chart data
				fetchTimeline(startDate, endDate, activeFilters).catch((err) => {
					console.error('Failed to load timeline:', err);
					return { timeline: [], timeline_format: 'day' };
				}),

				// Top pages
				fetchTopPages(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load pages:', err);
					return { top_pages: [] };
				}),

				// Entry and exit pages
				fetchEntryExitPages(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load entry/exit pages:', err);
					return { entry_pages: [], exit_pages: [] };
				}),

				// Top countries
				fetchTopCountries(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load countries:', err);
					return [];
				}),

				// Top sources
				fetchTopSources(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load sources:', err);
					return [];
				}),

				// Top events
				fetchTopEvents(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load events:', err);
					return [];
				}),

				// Browsers, devices, OS
				fetchBrowsersDevicesOS(startDate, endDate, 10, activeFilters).catch((err) => {
					console.error('Failed to load devices:', err);
					return { browsers: [], devices: [], os: [] };
				}),

				// Online users
				fetchOnlineUsers(5).catch((err) => {
					console.error('Failed to load online users:', err);
					return { online_users: 0, active_sessions: 0 };
				}),

				// Comparison stats (previous period)
				fetchTopStats(
					format(
						new Date(
							new Date(startDate).getTime() -
								(new Date(endDate).getTime() - new Date(startDate).getTime())
						),
						'yyyy-MM-dd'
					),
					format(subDays(new Date(startDate), 1), 'yyyy-MM-dd'),
					activeFilters
				).catch((err) => {
					console.error('Failed to load comparison stats:', err);
					return null;
				})
			]);

			// Update all state
			if (topStatsData) {
				stats = topStatsData;
			}

			if (timelineData) {
				timeline = timelineData;
			}

			if (pagesData) {
				topPages = pagesData;
			}

			if (entryExitData) {
				entryExitPages = entryExitData;
			}

			if (countriesData) {
				topCountries = countriesData;
			}

			if (sourcesData) {
				topSources = sourcesData;
			}

			if (eventsData) {
				topEvents = eventsData;
			}

			if (devicesData) {
				browsersDevicesOS = devicesData;
			}

			if (onlineData) {
				onlineData = onlineData;
			}

			if (comparisonData) {
				comparisonStats = comparisonData;
			}

			lastRefresh = new Date();
			updateURLParams();

			// All loading complete
			statsLoading = false;
			propertiesLoading = false;
			onlineUsersLoading = false;
		} catch (err: any) {
			error = err?.message || 'Failed to load stats';
			statsLoading = false;
			propertiesLoading = false;
			onlineUsersLoading = false;
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

	function handleRefreshIntervalChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		refreshIntervalTime = parseInt(target.value);
		setupAutoRefresh();
		updateURLParams();
	}

	function addFilter(type: keyof typeof activeFilters, value: string) {
		activeFilters[type] = value;
		updateURLParams();
		loadStats();
	}

	function removeFilter(type: keyof typeof activeFilters) {
		activeFilters[type] = null;
		updateURLParams();
		loadStats();
	}

	function clearAllFilters() {
		activeFilters = {
			source: null,
			country: null,
			browser: null,
			device: null,
			os: null,
			event: null,
			project: null,
			metric: null,
			propertyKey: null,
			propertyValue: null,
			botFilter: null,
			page: null
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
			// Failed to load projects - silently fail, not critical
		}
	}

	function handleDateChange(event: CustomEvent) {
		startDate = event.detail.startDate;
		endDate = event.detail.endDate;
		updateURLParams();
		loadStats();
	}

	// Helper to format trend
	function getTrendIcon(change: number) {
		if (change > 0) return TrendingUp;
		if (change < 0) return TrendingDown;
		return Minus;
	}

	function getTrendColor(change: number) {
		if (change > 0) return 'text-green-600';
		if (change < 0) return 'text-red-600';
		return 'text-gray-600';
	}

	// Handle metric card clicks for filtering timeline
	function handleMetricClick(metricType: string) {
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
	function isMetricSelected(metricType: string) {
		return activeFilters.metric === metricType;
	}

	// Handle property filter
	function addPropertyFilter(prop: { key: string; value: string }) {
		activeFilters.propertyKey = prop.key;
		activeFilters.propertyValue = prop.value;
		updateURLParams();
		loadStats();
	}

	// Comparison visibility state
	let showComparison = $state(true);
</script>

<div class="container mx-auto space-y-4 p-6">
	<!-- Header -->
	<div class="flex flex-wrap items-center justify-between gap-4">
		<div class="flex items-center gap-4">
			<h1 class="text-2xl font-bold">📊 {activeFilters.project || 'Siraaj'}</h1>
			{#if !loading}
				<Badge variant="secondary" class="gap-1">
					<span
						class="inline-block h-2 w-2 rounded-full {onlineData.online_users > 0
							? 'bg-green-500'
							: 'bg-gray-400'}"
					></span>
					{onlineData.online_users?.toLocaleString() || '0'} current visitors
				</Badge>
			{/if}
		</div>
	</div>

	<!-- Controls Row -->
	<div class="flex flex-wrap items-center gap-3">
		<!-- Date Range Selector -->
		<div class="flex items-center gap-2">
			<span class="text-sm font-medium">Period:</span>
			<select
				class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
				bind:value={dateRangePreset}
				onchange={(e: Event) => {
					const target = e.target as HTMLSelectElement;
					applyDateRangePreset(target.value);
				}}
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
				<span class="text-sm font-medium">Project:</span>
				<select
					class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
					value={activeFilters.project || ''}
					onchange={(e: Event) => {
						const target = e.target as HTMLSelectElement;
						if (target.value) {
							addFilter('project', target.value);
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

		<!-- Bot Filter -->
		<div class="flex items-center gap-2">
			<span class="text-sm font-medium">Traffic:</span>
			<select
				class="border-input bg-background focus-visible:ring-ring flex h-9 rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
				value={activeFilters.botFilter || ''}
				onchange={(e: Event) => {
					const target = e.target as HTMLSelectElement;
					if (target.value) {
						addFilter('botFilter', target.value);
					} else {
						removeFilter('botFilter');
					}
				}}
			>
				<option value="">All Traffic</option>
				<option value="human">👤 Human Only</option>
				<option value="bot">🤖 Bots Only</option>
			</select>
		</div>

		<!-- Auto-refresh controls -->
		<div class="ml-auto flex items-center gap-2">
			<span class="text-sm font-medium">Auto-refresh:</span>
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
			{#if activeFilters.device}
				<Badge variant="secondary" class="gap-1">
					Device: {activeFilters.device}
					<button onclick={() => removeFilter('device')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.os}
				<Badge variant="secondary" class="gap-1">
					OS: {activeFilters.os}
					<button onclick={() => removeFilter('os')} class="hover:text-destructive ml-1">
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
			{#if activeFilters.botFilter}
				<Badge variant="secondary" class="gap-1">
					Traffic: {activeFilters.botFilter === 'bot' ? '🤖 Bots Only' : '👤 Human Only'}
					<button onclick={() => removeFilter('botFilter')} class="hover:text-destructive ml-1">
						<X class="h-3 w-3" />
					</button>
				</Badge>
			{/if}
			{#if activeFilters.page}
				<Badge variant="secondary" class="gap-1">
					Page: {activeFilters.page}
					<button onclick={() => removeFilter('page')} class="hover:text-destructive ml-1">
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
		<!-- Timeline Chart -->
		<Card>
			<CardContent class="pb-2">
				<div class="mb-8 grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8">
					<!-- Unique Visitors -->
					<MetricCard
						label="Unique Visitors"
						currentValue={stats.unique_users || 0}
						previousValue={showComparison ? comparisonStats.unique_users || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						isSelected={isMetricSelected('users')}
						onclick={() => handleMetricClick('users')}
					/>

					<!-- Total Visits -->
					<MetricCard
						label="Total Visits"
						currentValue={stats.total_visits || 0}
						previousValue={showComparison ? comparisonStats.total_visits || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						isSelected={isMetricSelected('visits')}
						onclick={() => handleMetricClick('visits')}
					/>

					<!-- Total Pageviews -->
					<MetricCard
						label="Total Pageviews"
						currentValue={stats.page_views || 0}
						previousValue={showComparison ? comparisonStats.page_views || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						isSelected={isMetricSelected('page_views')}
						onclick={() => handleMetricClick('page_views')}
					/>

					<!-- Total Events -->
					<MetricCard
						label="Total Events"
						currentValue={stats.total_events || 0}
						previousValue={showComparison ? comparisonStats.total_events || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						isSelected={isMetricSelected('events')}
						onclick={() => handleMetricClick('events')}
					/>

					<!-- Views per Visit -->
					<MetricCard
						label="Views per Visit"
						currentValue={stats.total_visits > 0 ? stats.page_views / stats.total_visits : 0}
						previousValue={showComparison && comparisonStats.total_visits > 0
							? comparisonStats.page_views / comparisonStats.total_visits
							: null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						formatValue={(val) => (val ? val.toFixed(2) : '0.00')}
						isSelected={isMetricSelected('views_per_visit')}
						onclick={() => handleMetricClick('views_per_visit')}
					/>

					<!-- Bounce Rate -->
					<MetricCard
						label="Bounce Rate"
						currentValue={stats.bounce_rate || 0}
						previousValue={showComparison ? comparisonStats.bounce_rate || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						formatValue={(val) => (val ? val.toFixed(0) + '%' : '0%')}
						isNegativeBetter={true}
						isSelected={isMetricSelected('bounce_rate')}
						onclick={() => handleMetricClick('bounce_rate')}
					/>

					<!-- Visit Duration -->
					<MetricCard
						label="Visit Duration"
						currentValue={stats.avg_session_duration || 0}
						previousValue={showComparison ? comparisonStats.avg_session_duration || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						formatValue={(val) => {
							if (!val) return '0s';
							if (val < 60) return Math.floor(val) + 's';
							if (val < 3600) {
								const minutes = Math.floor(val / 60);
								const seconds = Math.floor(val % 60);
								return `${minutes}m ${seconds}s`;
							}
							const hours = Math.floor(val / 3600);
							const minutes = Math.floor((val % 3600) / 60);
							return `${hours}h ${minutes}m`;
						}}
						isSelected={isMetricSelected('visit_duration')}
						onclick={() => handleMetricClick('visit_duration')}
					/>

					<!-- Bot Percentage -->
					<MetricCard
						label="🤖 Bot Traffic"
						currentValue={stats.bot_percentage || 0}
						previousValue={showComparison ? comparisonStats.bot_percentage || 0 : null}
						currentPeriod={currentPeriodLabel()}
						previousPeriod={previousPeriodLabel()}
						formatValue={(val) => (val ? val.toFixed(0) + '%' : '0%')}
						isNegativeBetter={true}
						isSelected={activeFilters.botFilter === 'bot'}
						onclick={() => {
							if (activeFilters.botFilter === 'bot') {
								removeFilter('botFilter');
							} else {
								addFilter('botFilter', 'bot');
							}
						}}
					/>
				</div>

				{#if statsLoading}
					<div class="text-muted-foreground flex h-[300px] items-center justify-center">
						<div class="flex flex-col items-center gap-2">
							<div
								class="border-primary h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
							></div>
							<p class="text-sm">Loading...</p>
						</div>
					</div>
				{:else}
					<TimelineChart
						data={timeline.timeline || []}
						comparisonData={[]}
						format={timeline.timeline_format || 'day'}
						metric={activeFilters.metric || 'users'}
						bind:showComparison
					/>
				{/if}
			</CardContent>
		</Card>

		<!-- Data Grid -->
		<div class="grid gap-4 lg:grid-cols-2">
			<!-- Left Column -->
			<div class="space-y-4">
				<Card>
					<CardHeader class="pb-3">
						<div class="flex items-center justify-between">
							<CardTitle class="text-base">Pages</CardTitle>
							<div class="bg-muted flex gap-1 rounded-lg p-1">
								<button
									class="rounded px-3 py-1 text-xs font-medium transition-colors {pagesTab === 'all'
										? 'bg-background shadow-sm'
										: 'hover:bg-background/50'}"
									onclick={() => (pagesTab = 'all')}
								>
									All
								</button>
								<button
									class="rounded px-3 py-1 text-xs font-medium transition-colors {pagesTab ===
									'entry'
										? 'bg-background shadow-sm'
										: 'hover:bg-background/50'}"
									onclick={() => (pagesTab = 'entry')}
								>
									Entry
								</button>
								<button
									class="rounded px-3 py-1 text-xs font-medium transition-colors {pagesTab ===
									'exit'
										? 'bg-background shadow-sm'
										: 'hover:bg-background/50'}"
									onclick={() => (pagesTab = 'exit')}
								>
									Exit
								</button>
							</div>
						</div>
					</CardHeader>
					<CardContent>
						{#if statsLoading}
							<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
								<div class="flex flex-col items-center gap-2">
									<div
										class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
									></div>
									<p class="text-xs">Loading...</p>
								</div>
							</div>
						{:else if pagesTab === 'all'}
							<TopItemsList
								items={topPages.top_pages || []}
								labelKey="url"
								maxItems={10}
								valueKey="count"
								showMoreTitle="All Pages"
								onclick={(item: any) => addFilter('page', item.url)}
							/>
						{:else if pagesTab === 'entry'}
							<TopItemsList
								items={entryExitPages.entry_pages || []}
								labelKey="url"
								maxItems={10}
								valueKey="count"
								showMoreTitle="All Entry Pages"
								onclick={(item: any) => addFilter('page', item.url)}
							/>
						{:else if pagesTab === 'exit'}
							<TopItemsList
								items={entryExitPages.exit_pages || []}
								labelKey="url"
								maxItems={10}
								valueKey="count"
								showMoreTitle="All Exit Pages"
								onclick={(item: any) => addFilter('page', item.url)}
							/>
						{/if}
					</CardContent>
				</Card>
			</div>
			<Card>
				<CardHeader class="pb-3">
					<CardTitle class="text-base">Locations</CardTitle>
				</CardHeader>
				<CardContent>
					<CountriesPanel
						countries={topCountries || []}
						onclick={(item: any) => addFilter('country', item.name)}
						loading={statsLoading}
					/>
				</CardContent>
			</Card>
		</div>

		<div class="grid gap-4 lg:grid-cols-3">
			<Card>
				<CardHeader class="pb-3">
					<CardTitle class="text-base">Top Events</CardTitle>
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
							items={topEvents || []}
							labelKey="name"
							valueKey="count"
							maxItems={8}
							type="event"
							showMoreTitle="All Events"
							onclick={(item: any) => addFilter('event', item.name)}
						/>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader class="pb-3">
					<CardTitle class="text-base">Top Sources</CardTitle>
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
							items={topSources || []}
							labelKey="name"
							valueKey="count"
							maxItems={8}
							type="source"
							showMoreTitle="All Sources"
							onclick={(item: any) => addFilter('source', item.name)}
						/>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader class="pb-3">
					<CardTitle class="text-base">Devices</CardTitle>
				</CardHeader>
				<CardContent>
					<BrowserPanel
						browsers={browsersDevicesOS.browsers || []}
						devices={browsersDevicesOS.devices || []}
						operatingSystems={browsersDevicesOS.os || []}
						onBrowserClick={(item: any) => addFilter('browser', item.name)}
						onDeviceClick={(item: any) => addFilter('device', item.name)}
						onOsClick={(item: any) => addFilter('os', item.name)}
						loading={statsLoading}
					/>
				</CardContent>
			</Card>
		</div>
	{/if}
</div>
