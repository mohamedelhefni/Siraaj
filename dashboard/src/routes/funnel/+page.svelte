<script lang="ts">
	import { onMount } from 'svelte';
	import { format, subDays } from 'date-fns';
	import { fetchFunnelAnalysis, fetchStats, fetchProjects } from '$lib/api';
	import { Play, Save, FolderOpen, Download, Share2 } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import FunnelBuilder from '$lib/components/FunnelBuilder.svelte';
	import FunnelChart from '$lib/components/FunnelChart.svelte';

	interface FunnelStep {
		name: string;
		event_name: string;
		url?: string;
		filters?: Record<string, string>;
		expanded?: boolean;
	}

	let steps = $state<FunnelStep[]>([]);
	let funnelResult = $state<any>(null);
	let loading = $state(false);
	let error = $state<string | null>(null);

	// Date range
	let startDate = $state(format(subDays(new Date(), 7), 'yyyy-MM-dd'));
	let endDate = $state(format(new Date(), 'yyyy-MM-dd'));

	// Filters
	let selectedProject = $state('');
	let selectedCountry = $state('');
	let selectedBrowser = $state('');
	let selectedDevice = $state('');
	let selectedOs = $state('');
	let botFilter = $state('');

	// Available options for dropdowns
	let projects = $state<string[]>([]);
	let availableEvents = $state<string[]>([]);
	let availablePages = $state<string[]>([]);

	onMount(() => {
		loadProjects();
		loadAvailableOptions();
		loadFromURL();
	});

	async function loadProjects() {
		try {
			projects = await fetchProjects();
		} catch (err) {
			console.error('Failed to load projects:', err);
		}
	}

	async function loadAvailableOptions() {
		try {
			const stats = await fetchStats(
				format(subDays(new Date(), 30), 'yyyy-MM-dd'),
				format(new Date(), 'yyyy-MM-dd'),
				100,
				{}
			);

			// Extract unique events
			if (stats.top_events) {
				availableEvents = stats.top_events.map((e: any) => e.name);
			}

			// Extract unique pages
			if (stats.top_pages) {
				availablePages = stats.top_pages.map((p: any) => p.url);
			}
		} catch (err) {
			console.error('Failed to load available options:', err);
		}
	}

	// Watch for steps changes and update URL
	$effect(() => {
		if (steps.length > 0) {
			updateURL();
		}
	});

	async function runAnalysis() {
		if (steps.length === 0) {
			error = 'Please add at least one funnel step';
			return;
		}

		// Validate steps
		const invalidStep = steps.find((s) => !s.event_name);
		if (invalidStep) {
			error = 'All steps must have an event name';
			return;
		}

		loading = true;
		error = null;
		funnelResult = null;

		try {
			const filters: Record<string, string> = {};
			if (selectedProject) filters.project = selectedProject;
			if (selectedCountry) filters.country = selectedCountry;
			if (selectedBrowser) filters.browser = selectedBrowser;
			if (selectedDevice) filters.device = selectedDevice;
			if (selectedOs) filters.os = selectedOs;
			if (botFilter) filters.botFilter = botFilter;

			const request = {
				steps: steps.map((s) => ({
					name: s.name,
					event_name: s.event_name,
					url: s.url || '',
					filters: s.filters || {}
				})),
				start_date: startDate,
				end_date: endDate,
				filters
			};

			console.log('Running funnel analysis:', request);

			const result = await fetchFunnelAnalysis(request);
			funnelResult = result;

			updateURL();
		} catch (err: any) {
			error = err?.message || 'Failed to analyze funnel';
			console.error('Funnel analysis error:', err);
		} finally {
			loading = false;
		}
	}

	function updateURL() {
		if (typeof window === 'undefined') return;

		const params = new URLSearchParams();
		params.set('start', startDate);
		params.set('end', endDate);

		if (selectedProject) params.set('project', selectedProject);
		if (selectedCountry) params.set('country', selectedCountry);
		if (selectedBrowser) params.set('browser', selectedBrowser);
		if (selectedDevice) params.set('device', selectedDevice);
		if (selectedOs) params.set('os', selectedOs);
		if (botFilter) params.set('botFilter', botFilter);

		if (steps.length > 0) {
			params.set(
				'steps',
				JSON.stringify(
					steps.map((s) => ({
						name: s.name,
						event_name: s.event_name,
						url: s.url || ''
					}))
				)
			);
		}

		const newURL = `${window.location.pathname}?${params.toString()}`;
		window.history.replaceState({}, '', newURL);
	}

	function loadFromURL() {
		if (typeof window === 'undefined') return;

		const params = new URLSearchParams(window.location.search);

		const start = params.get('start');
		const end = params.get('end');
		if (start) startDate = start;
		if (end) endDate = end;

		const project = params.get('project');
		const country = params.get('country');
		const browser = params.get('browser');
		const device = params.get('device');
		const os = params.get('os');
		const bot = params.get('botFilter');

		if (project) selectedProject = project;
		if (country) selectedCountry = country;
		if (browser) selectedBrowser = browser;
		if (device) selectedDevice = device;
		if (os) selectedOs = os;
		if (bot) botFilter = bot;

		const stepsParam = params.get('steps');
		if (stepsParam) {
			try {
				const parsedSteps = JSON.parse(stepsParam);
				steps = parsedSteps.map((s: any) => ({
					...s,
					expanded: false
				}));
			} catch (err) {
				console.error('Failed to parse steps from URL:', err);
			}
		}
	}

	function exportResults() {
		if (!funnelResult) return;

		const dataStr = JSON.stringify(funnelResult, null, 2);
		const dataBlob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(dataBlob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `funnel-analysis-${format(new Date(), 'yyyy-MM-dd-HHmmss')}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}

	function shareResults() {
		if (typeof window === 'undefined') return;

		const url = window.location.href;
		navigator.clipboard
			.writeText(url)
			.then(() => {
				alert('Link copied to clipboard!');
			})
			.catch((err) => {
				console.error('Failed to copy link:', err);
			});
	}
</script>

<div class="container mx-auto space-y-6 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">🎯 Funnel Analysis</h1>
			<p class="text-muted-foreground mt-1">
				Track user journey through multi-step conversion funnels
			</p>
		</div>

		{#if funnelResult}
			<div class="flex items-center gap-2">
				<Button variant="outline" size="sm" onclick={exportResults}>
					<Download class="mr-1 h-4 w-4" />
					Export
				</Button>
				<Button variant="outline" size="sm" onclick={shareResults}>
					<Share2 class="mr-1 h-4 w-4" />
					Share
				</Button>
			</div>
		{/if}
	</div>

	<div class="grid gap-6 lg:grid-cols-3">
		<!-- Left Panel: Configuration -->
		<div class="space-y-4 lg:col-span-1">
			<!-- Date Range & Filters -->
			<Card>
				<CardHeader>
					<CardTitle>Configuration</CardTitle>
					<CardDescription>Set date range and filters for your funnel</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<!-- Date Range -->
					<div>
						<label class="mb-2 block text-sm font-medium">Date Range</label>
						<div class="space-y-2">
							<input
								type="date"
								bind:value={startDate}
								class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
							/>
							<input
								type="date"
								bind:value={endDate}
								class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
							/>
						</div>
					</div>

					<!-- Global Filters -->
					<div class="space-y-3">
						<div class="text-sm font-medium">Global Filters</div>

						{#if projects.length > 0}
							<div>
								<label class="text-muted-foreground mb-1 block text-xs">Project</label>
								<select
									bind:value={selectedProject}
									class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
								>
									<option value="">All Projects</option>
									{#each projects as project}
										<option value={project}>{project}</option>
									{/each}
								</select>
							</div>
						{/if}

						<div>
							<label class="text-muted-foreground mb-1 block text-xs">Traffic Type</label>
							<select
								bind:value={botFilter}
								class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
							>
								<option value="">All Traffic</option>
								<option value="human">👤 Human Only</option>
								<option value="bot">🤖 Bots Only</option>
							</select>
						</div>
					</div>

					<!-- Run Button -->
					<Button onclick={runAnalysis} disabled={loading || steps.length === 0} class="w-full">
						{#if loading}
							<div
								class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
							></div>
							Analyzing...
						{:else}
							<Play class="mr-2 h-4 w-4" />
							Run Analysis
						{/if}
					</Button>

					{#if error}
						<div class="bg-destructive/10 text-destructive rounded-lg p-3 text-sm">
							{error}
						</div>
					{/if}
				</CardContent>
			</Card>

			<!-- Quick Stats -->
			{#if funnelResult}
				<Card>
					<CardHeader>
						<CardTitle class="text-base">Summary</CardTitle>
					</CardHeader>
					<CardContent class="space-y-2">
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Time Range:</span>
							<span class="font-medium">{funnelResult.time_range}</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Total Steps:</span>
							<span class="font-medium">{funnelResult.steps.length}</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Users Entered:</span>
							<span class="font-medium">{funnelResult.total_users.toLocaleString()}</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Users Completed:</span>
							<span class="font-medium">{funnelResult.completed_users.toLocaleString()}</span>
						</div>
					</CardContent>
				</Card>
			{/if}
		</div>

		<!-- Right Panel: Funnel Builder & Results -->
		<div class="space-y-6 lg:col-span-2">
			<!-- Funnel Builder -->
			<Card>
				<CardHeader>
					<CardTitle>Funnel Steps</CardTitle>
					<CardDescription>
						Define the sequence of events that make up your conversion funnel
					</CardDescription>
				</CardHeader>
				<CardContent>
					<FunnelBuilder bind:steps {availableEvents} {availablePages} />
				</CardContent>
			</Card>

			<!-- Results -->
			<FunnelChart result={funnelResult} {loading} />
		</div>
	</div>
</div>
