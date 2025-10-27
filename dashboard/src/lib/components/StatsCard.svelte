<script>
	import { Activity, Users, TrendingUp, Globe, BarChart3, Eye } from 'lucide-svelte';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { formatCompactNumber } from '$lib/utils/formatters.js';

	let { title, value, icon, description } = $props();

	const icons = {
		activity: Activity,
		users: Users,
		'trending-up': TrendingUp,
		globe: Globe,
		'bar-chart': BarChart3,
		eye: Eye
	};

	const IconComponent = icons[icon] || Activity;
	
	// Auto-format numeric values
	const formattedValue = $derived(() => {
		if (typeof value === 'number') {
			return formatCompactNumber(value);
		}
		return value;
	});
</script>

<Card>
	<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
		<CardTitle class="text-sm font-medium">{title}</CardTitle>
		<IconComponent class="text-muted-foreground h-4 w-4" />
	</CardHeader>
	<CardContent>
		<div class="text-2xl font-bold">{formattedValue()}</div>
		<p class="text-muted-foreground mt-1 text-xs">{description}</p>
	</CardContent>
</Card>
