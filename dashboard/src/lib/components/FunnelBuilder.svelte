<script lang="ts">
	import { Plus, X, Trash2, GripVertical, ChevronDown, ChevronUp } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';

	interface FunnelStep {
		name: string;
		event_name: string;
		url?: string;
		filters?: Record<string, string>;
		expanded?: boolean;
	}

	interface Props {
		steps: FunnelStep[];
		availableEvents?: string[];
		availablePages?: string[];
	}

	let { steps = $bindable([]), availableEvents = [], availablePages = [] }: Props = $props();

	function addStep() {
		const newStep: FunnelStep = {
			name: `Step ${steps.length + 1}`,
			event_name: '',
			url: '',
			filters: {},
			expanded: true
		};
		steps = [...steps, newStep];
	}

	function removeStep(index: number) {
		steps = steps.filter((_, i) => i !== index);
	}

	function updateStep(index: number, field: keyof FunnelStep, value: any) {
		steps = steps.map((step, i) => (i === index ? { ...step, [field]: value } : step));
	}

	function toggleExpanded(index: number) {
		steps = steps.map((step, i) => (i === index ? { ...step, expanded: !step.expanded } : step));
	}

	function moveStep(index: number, direction: 'up' | 'down') {
		if (direction === 'up' && index > 0) {
			const newSteps = [...steps];
			[newSteps[index - 1], newSteps[index]] = [newSteps[index], newSteps[index - 1]];
			steps = newSteps;
		} else if (direction === 'down' && index < steps.length - 1) {
			const newSteps = [...steps];
			[newSteps[index], newSteps[index + 1]] = [newSteps[index + 1], newSteps[index]];
			steps = newSteps;
		}
	}

	function clearAll() {
		if (confirm('Are you sure you want to clear all funnel steps?')) {
			steps = [];
		}
	}

	// Quick templates
	function loadTemplate(template: 'signup' | 'purchase' | 'engagement') {
		let templateSteps: FunnelStep[] = [];

		switch (template) {
			case 'signup':
				templateSteps = [
					{ name: 'Landing Page', event_name: 'page_view', url: '/', expanded: false },
					{ name: 'Signup Page', event_name: 'page_view', url: '/signup', expanded: false },
					{ name: 'Signup Complete', event_name: 'signup', expanded: false }
				];
				break;
			case 'purchase':
				templateSteps = [
					{ name: 'Product View', event_name: 'product_view', expanded: false },
					{ name: 'Add to Cart', event_name: 'add_to_cart', expanded: false },
					{ name: 'Checkout', event_name: 'checkout', expanded: false },
					{ name: 'Purchase', event_name: 'purchase', expanded: false }
				];
				break;
			case 'engagement':
				templateSteps = [
					{ name: 'First Visit', event_name: 'page_view', expanded: false },
					{ name: 'Feature Click', event_name: 'button_clicked', expanded: false },
					{ name: 'Return Visit', event_name: 'page_view', expanded: false }
				];
				break;
		}

		steps = templateSteps;
	}
</script>

<div class="space-y-4">
	<!-- Header with quick actions -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-2">
			<h3 class="text-lg font-semibold">Funnel Steps</h3>
			{#if steps.length > 0}
				<Badge variant="secondary">{steps.length} step{steps.length !== 1 ? 's' : ''}</Badge>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<!-- Quick Templates -->
			<div class="flex gap-1">
				<Button variant="outline" size="sm" onclick={() => loadTemplate('signup')}>
					📝 Signup
				</Button>
				<Button variant="outline" size="sm" onclick={() => loadTemplate('purchase')}>
					🛒 Purchase
				</Button>
				<Button variant="outline" size="sm" onclick={() => loadTemplate('engagement')}>
					📊 Engagement
				</Button>
			</div>
			{#if steps.length > 0}
				<Button variant="ghost" size="sm" onclick={clearAll}>
					<Trash2 class="mr-1 h-4 w-4" />
					Clear All
				</Button>
			{/if}
			<Button onclick={addStep} size="sm">
				<Plus class="mr-1 h-4 w-4" />
				Add Step
			</Button>
		</div>
	</div>

	<!-- Steps List -->
	{#if steps.length === 0}
		<Card class="border-dashed">
			<CardContent class="flex flex-col items-center justify-center py-12">
				<div class="text-muted-foreground mb-4 text-center">
					<p class="mb-2 text-lg font-medium">No funnel steps defined</p>
					<p class="text-sm">Add steps to analyze user journey through your funnel</p>
				</div>
				<Button onclick={addStep}>
					<Plus class="mr-2 h-4 w-4" />
					Add First Step
				</Button>
			</CardContent>
		</Card>
	{:else}
		<div class="space-y-3">
			{#each steps as step, index}
				<Card class="relative">
					<CardContent class="p-4">
						<!-- Step Header -->
						<div class="flex items-start gap-3">
							<!-- Step Number -->
							<div
								class="bg-primary text-primary-foreground flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-sm font-bold"
							>
								{index + 1}
							</div>

							<!-- Step Content -->
							<div class="flex-1 space-y-3">
								<!-- Collapsed View -->
								<div class="flex items-center justify-between">
									<div class="flex-1">
										<input
											type="text"
											value={step.name}
											oninput={(e) =>
												updateStep(index, 'name', (e.target as HTMLInputElement).value)}
											placeholder="Step name"
											class="w-full border-none bg-transparent p-0 text-base font-medium focus:outline-none focus:ring-0"
										/>
										<div
											class="text-muted-foreground mt-1 flex flex-wrap items-center gap-2 text-sm"
										>
											{#if step.event_name}
												<Badge variant="outline" class="text-xs">
													Event: {step.event_name}
												</Badge>
											{/if}
											{#if step.url}
												<Badge variant="outline" class="text-xs">
													URL: {step.url}
												</Badge>
											{/if}
											{#if !step.event_name && !step.url}
												<span class="text-destructive text-xs">⚠️ Configure step details</span>
											{/if}
										</div>
									</div>

									<div class="flex items-center gap-1">
										<!-- Expand/Collapse -->
										<Button variant="ghost" size="sm" onclick={() => toggleExpanded(index)}>
											{#if step.expanded}
												<ChevronUp class="h-4 w-4" />
											{:else}
												<ChevronDown class="h-4 w-4" />
											{/if}
										</Button>

										<!-- Move buttons -->
										{#if index > 0}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => moveStep(index, 'up')}
												title="Move up"
											>
												↑
											</Button>
										{/if}
										{#if index < steps.length - 1}
											<Button
												variant="ghost"
												size="sm"
												onclick={() => moveStep(index, 'down')}
												title="Move down"
											>
												↓
											</Button>
										{/if}

										<!-- Delete button -->
										<Button variant="ghost" size="sm" onclick={() => removeStep(index)}>
											<X class="text-destructive h-4 w-4" />
										</Button>
									</div>
								</div>

								<!-- Expanded View -->
								{#if step.expanded}
									<div class="space-y-3 border-t pt-3">
										<!-- Event Name -->
										<div>
											<label class="text-muted-foreground mb-1.5 block text-xs font-medium">
												Event Name *
											</label>
											{#if availableEvents.length > 0}
												<select
													value={step.event_name}
													onchange={(e) =>
														updateStep(index, 'event_name', (e.target as HTMLSelectElement).value)}
													class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
												>
													<option value="">Select event...</option>
													{#each availableEvents as event}
														<option value={event}>{event}</option>
													{/each}
												</select>
											{:else}
												<input
													type="text"
													value={step.event_name}
													oninput={(e) =>
														updateStep(index, 'event_name', (e.target as HTMLInputElement).value)}
													placeholder="e.g., page_view, button_clicked"
													class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
												/>
											{/if}
										</div>

										<!-- URL (Optional) -->
										<div>
											<label class="text-muted-foreground mb-1.5 block text-xs font-medium">
												URL Pattern (optional)
											</label>
											{#if availablePages.length > 0}
												<select
													value={step.url || ''}
													onchange={(e) =>
														updateStep(index, 'url', (e.target as HTMLSelectElement).value)}
													class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
												>
													<option value="">Any URL</option>
													{#each availablePages as page}
														<option value={page}>{page}</option>
													{/each}
												</select>
											{:else}
												<input
													type="text"
													value={step.url || ''}
													oninput={(e) =>
														updateStep(index, 'url', (e.target as HTMLInputElement).value)}
													placeholder="e.g., /signup, /checkout"
													class="border-input bg-background focus-visible:ring-ring flex h-9 w-full rounded-md border px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1"
												/>
											{/if}
											<p class="text-muted-foreground mt-1 text-xs">
												Match events from a specific URL
											</p>
										</div>
									</div>
								{/if}
							</div>
						</div>

						<!-- Connection Line to Next Step -->
						{#if index < steps.length - 1}
							<div class="text-muted-foreground ml-4 mt-2 flex items-center gap-2">
								<div class="bg-border ml-3.5 h-8 w-px"></div>
								<span class="text-xs">↓ Then</span>
							</div>
						{/if}
					</CardContent>
				</Card>
			{/each}
			<Button onclick={addStep} size="sm" class="w-full cursor-pointer">
				<Plus class="mr-1 h-4 w-4" />
				Add Step
			</Button>
		</div>
	{/if}

	{#if steps.length > 0}
		<div class="bg-muted/50 rounded-lg p-3 text-sm">
			<p class="text-muted-foreground">
				💡 <strong>Tip:</strong> Funnel analysis tracks users who complete steps in order. Each step
				must occur after the previous one for the same user.
			</p>
		</div>
	{/if}
</div>
