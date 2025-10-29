<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/stores';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen">
	<!-- Navigation -->
	<nav class="bg-background border-b">
		<div class="container mx-auto flex items-center gap-6 px-6 py-3">
			<a href="/dashboard" class="flex items-center gap-2 text-lg font-bold"> 🏮 Siraaj </a>
			<div class="flex gap-1">
				<a
					href="/dashboard"
					class="hover:bg-accent rounded-lg px-3 py-2 text-sm font-medium transition-colors {$page
						.url.pathname === '/dashboard' || $page.url.pathname === '/'
						? 'bg-accent'
						: ''}"
				>
					📊 Dashboard
				</a>
				<a
					href="/dashboard/channels"
					class="hover:bg-accent rounded-lg px-3 py-2 text-sm font-medium transition-colors {$page
						.url.pathname === '/channels'
						? 'bg-accent'
						: ''}"
				>
					📈 Channels
				</a>
				<a
					href="/dashboard/funnel"
					class="hover:bg-accent rounded-lg px-3 py-2 text-sm font-medium transition-colors {$page
						.url.pathname === '/dashboard/funnel'
						? 'bg-accent'
						: ''}"
				>
					🎯 Funnels
				</a>
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	{@render children?.()}
</div>
