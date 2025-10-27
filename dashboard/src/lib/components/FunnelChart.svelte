<script lang="ts">
	import { TrendingDown, TrendingUp, Clock, Users, Activity } from 'lucide-svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCompactNumber } from '$lib/utils/formatters.js';

	interface FunnelStepResult {
		step: {
			name: string;
			event_name: string;
			url?: string;
		};
		user_count: number;
		session_count: number;
		event_count: number;
		conversion_rate: number;
		overall_rate: number;
		dropoff_rate: number;
		avg_time_to_next?: number;
		median_time_to_next?: number;
	}

	interface Props {
		result: {
			steps: FunnelStepResult[];
			total_users: number;
			completed_users: number;
			completion_rate: number;
			avg_completion?: number;
			time_range: string;
		} | null;
		loading?: boolean;
	}

	let { result, loading = false }: Props = $props();

	function formatTime(seconds: number | undefined): string {
		if (!seconds || seconds === 0) return 'N/A';

		if (seconds < 60) {
			return `${Math.round(seconds)}s`;
		} else if (seconds < 3600) {
			const mins = Math.floor(seconds / 60);
			const secs = Math.round(seconds % 60);
			return `${mins}m ${secs}s`;
		} else {
			const hours = Math.floor(seconds / 3600);
			const mins = Math.round((seconds % 3600) / 60);
			return `${hours}h ${mins}m`;
		}
	}

	function getBarWidth(userCount: number, maxUsers: number): number {
		if (maxUsers === 0) return 0;
		return (userCount / maxUsers) * 100;
	}

	function getConversionColor(rate: number): string {
		if (rate >= 75) return 'text-green-600';
		if (rate >= 50) return 'text-yellow-600';
		if (rate >= 25) return 'text-orange-600';
		return 'text-red-600';
	}

	function getConversionBg(rate: number): string {
		if (rate >= 75) return 'bg-green-500';
		if (rate >= 50) return 'bg-yellow-500';
		if (rate >= 25) return 'bg-orange-500';
		return 'bg-red-500';
	}

	$effect(() => {
		if (result) {
			console.log('Funnel Result:', result);
		}
	});
</script>

{#if loading}
	<div class="flex min-h-[400px] items-center justify-center">
		<div class="flex flex-col items-center gap-3">
			<div
				class="border-primary h-10 w-10 animate-spin rounded-full border-4 border-t-transparent"
			></div>
			<p class="text-muted-foreground text-sm">Analyzing funnel...</p>
		</div>
	</div>
{:else if !result}
	<div
		class="bg-muted/30 flex min-h-[400px] flex-col items-center justify-center rounded-lg border-2 border-dashed"
	>
		<div class="text-muted-foreground text-center">
			<Activity class="mx-auto mb-3 h-12 w-12 opacity-50" />
			<p class="mb-2 text-lg font-medium">No funnel data</p>
			<p class="text-sm">Configure your funnel steps and run the analysis</p>
		</div>
	</div>
{:else}
	<div class="space-y-6">
		<!-- Summary Stats -->
		<div class="grid gap-4 md:grid-cols-4">
			<Card>
				<CardContent class="p-4">
					<div class="text-muted-foreground mb-1 text-xs font-medium uppercase">Total Entered</div>
					<div class="text-2xl font-bold">{formatCompactNumber(result.total_users)}</div>
					<div class="text-muted-foreground mt-1 text-xs">users started the funnel</div>
				</CardContent>
			</Card>

			<Card>
				<CardContent class="p-4">
					<div class="text-muted-foreground mb-1 text-xs font-medium uppercase">Completed</div>
					<div class="text-2xl font-bold">{formatCompactNumber(result.completed_users)}</div>
					<div class="text-muted-foreground mt-1 text-xs">users completed all steps</div>
				</CardContent>
			</Card>

			<Card>
				<CardContent class="p-4">
					<div class="text-muted-foreground mb-1 text-xs font-medium uppercase">
						Completion Rate
					</div>
					<div class="flex items-baseline gap-2">
						<div class="text-2xl font-bold {getConversionColor(result.completion_rate)}">
							{result.completion_rate.toFixed(1)}%
						</div>
					</div>
					<div class="text-muted-foreground mt-1 text-xs">end-to-end conversion</div>
				</CardContent>
			</Card>

			<Card>
				<CardContent class="p-4">
					<div class="text-muted-foreground mb-1 text-xs font-medium uppercase">
						Avg. Completion Time
					</div>
					<div class="text-2xl font-bold">{formatTime(result.avg_completion)}</div>
					<div class="text-muted-foreground mt-1 text-xs">to complete funnel</div>
				</CardContent>
			</Card>
		</div>

		<!-- Funnel Visualization -->
		<Card>
			<CardHeader>
				<CardTitle>Funnel Steps</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-4">
					{#each result.steps as stepResult, index}
						{@const maxUsers = result.total_users}
						{@const barWidth = getBarWidth(stepResult.user_count, maxUsers)}

						<div class="space-y-2">
							<!-- Step Header -->
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-3">
									<div
										class="bg-primary text-primary-foreground flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-sm font-bold"
									>
										{index + 1}
									</div>
									<div>
										<div class="font-medium">{stepResult.step.name}</div>
										<div class="text-muted-foreground flex flex-wrap items-center gap-2 text-xs">
											<span>Event: {stepResult.step.event_name}</span>
											{#if stepResult.step.url}
												<span>• URL: {stepResult.step.url}</span>
											{/if}
										</div>
									</div>
								</div>

								<div class="text-right">
									<div class="text-lg font-bold">
										{formatCompactNumber(stepResult.user_count)}
									</div>
									<div class="text-muted-foreground text-xs">users</div>
								</div>
							</div>
							<!-- Funnel Bar -->
							<div class="bg-muted relative h-16 overflow-hidden rounded-lg">
								<div
									class="{getConversionBg(
										stepResult.conversion_rate
									)} absolute inset-y-0 left-0 transition-all duration-500"
									style="width: {barWidth}%"
								>
									<div class="flex h-full items-center justify-between px-4 text-white">
										<div class="flex items-center gap-4">
											<div class="text-sm font-semibold">
												{stepResult.overall_rate.toFixed(1)}% of total
											</div>
											{#if index > 0}
												<div class="border-l border-white/30 pl-4 text-sm">
													{stepResult.conversion_rate.toFixed(1)}% from previous
												</div>
											{/if}
										</div>
									</div>
								</div>
							</div>

							<!-- Step Stats Grid -->
							<div class="grid gap-3 sm:grid-cols-3">
								<div class="bg-muted/50 rounded-lg p-3">
									<div class="text-muted-foreground mb-1 flex items-center gap-1 text-xs">
										<Users class="h-3 w-3" />
										Sessions
									</div>
									<div class="font-semibold">{formatCompactNumber(stepResult.session_count)}</div>
								</div>

								<div class="bg-muted/50 rounded-lg p-3">
									<div class="text-muted-foreground mb-1 flex items-center gap-1 text-xs">
										<Activity class="h-3 w-3" />
										Events
									</div>
									<div class="font-semibold">{formatCompactNumber(stepResult.event_count)}</div>
								</div>
								{#if stepResult.avg_time_to_next && stepResult.avg_time_to_next > 0}
									<div class="bg-muted/50 rounded-lg p-3">
										<div class="text-muted-foreground mb-1 flex items-center gap-1 text-xs">
											<Clock class="h-3 w-3" />
											Time to Next
										</div>
										<div class="font-semibold">
											{formatTime(stepResult.avg_time_to_next)}
											<span class="text-muted-foreground text-xs">avg</span>
										</div>
									</div>
								{/if}
							</div>

							<!-- Dropoff Indicator -->
							{#if index < result.steps.length - 1 && stepResult.dropoff_rate > 0}
								<div class="flex items-center gap-2 pl-10">
									<TrendingDown class="text-destructive h-4 w-4" />
									<span class="text-destructive text-sm font-medium">
										{stepResult.dropoff_rate.toFixed(1)}% drop-off
									</span>
									<span class="text-muted-foreground text-xs">
										({formatCompactNumber(
											result.steps[index + 1].user_count - stepResult.user_count
										)}
										users left)
									</span>
								</div>
							{/if}
							<!-- Separator -->
							{#if index < result.steps.length - 1}
								<div class="border-t"></div>
							{/if}
						</div>
					{/each}
				</div>
			</CardContent>
		</Card>

		<!-- Insights -->
		<Card>
			<CardHeader>
				<CardTitle>Key Insights</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-3">
					{#if result.steps.length > 0}
						{@const biggestDropoff = result.steps
							.slice(0, -1)
							.reduce(
								(max, step, idx) =>
									step.dropoff_rate > (result.steps[max]?.dropoff_rate || 0) ? idx : max,
								0
							)}

						{#if result.steps[biggestDropoff]?.dropoff_rate > 0}
							<div class="bg-destructive/10 border-destructive/20 rounded-lg border p-4">
								<div class="mb-1 flex items-start gap-2">
									<TrendingDown class="text-destructive mt-0.5 h-5 w-5" />
									<div>
										<div class="font-semibold">Biggest Drop-off Point</div>
										<p class="text-muted-foreground text-sm">
											{result.steps[biggestDropoff].dropoff_rate.toFixed(1)}% of users drop off
											after
											<strong>"{result.steps[biggestDropoff].step.name}"</strong>
										</p>
									</div>
								</div>
							</div>
						{/if}

						{@const bestConversion = result.steps
							.slice(1)
							.reduce(
								(max, step, idx) =>
									step.conversion_rate > (result.steps[max + 1]?.conversion_rate || 0)
										? idx + 1
										: max + 1,
								1
							)}

						{#if result.steps.length > 1 && result.steps[bestConversion]}
							<div
								class="rounded-lg border border-green-200 bg-green-50 p-4 dark:border-green-900 dark:bg-green-950/20"
							>
								<div class="mb-1 flex items-start gap-2">
									<TrendingUp class="mt-0.5 h-5 w-5 text-green-600 dark:text-green-500" />
									<div>
										<div class="font-semibold">Best Performing Step</div>
										<p class="text-muted-foreground text-sm">
											{result.steps[bestConversion].conversion_rate.toFixed(1)}% conversion rate
											from previous step at
											<strong>"{result.steps[bestConversion].step.name}"</strong>
										</p>
									</div>
								</div>
							</div>
						{/if}

						{#if result.completion_rate < 10}
							<div
								class="rounded-lg border border-yellow-200 bg-yellow-50 p-4 dark:border-yellow-900 dark:bg-yellow-950/20"
							>
								<div class="mb-1 flex items-start gap-2">
									<span class="text-xl">⚠️</span>
									<div>
										<div class="font-semibold">Low Completion Rate</div>
										<p class="text-muted-foreground text-sm">
											Only {result.completion_rate.toFixed(1)}% of users complete the entire funnel.
											Consider simplifying the journey or reducing friction.
										</p>
									</div>
								</div>
							</div>
						{/if}
					{/if}
				</div>
			</CardContent>
		</Card>
	</div>
{/if}
