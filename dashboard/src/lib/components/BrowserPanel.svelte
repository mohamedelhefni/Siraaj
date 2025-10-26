<script lang="ts">
	import { Monitor, Smartphone, Laptop } from 'lucide-svelte';
	import TopItemsList from './TopItemsList.svelte';
	import { Button } from '$lib/components/ui/button';

	let {
		browsers = [],
		devices = [],
		operatingSystems = [],
		onBrowserClick = null,
		onDeviceClick = null,
		onOsClick = null,
		loading = false
	} = $props();

	let activeTab = $state('browsers'); // Default to 'browsers' view

	function handleBrowserClick(item: { name: string; count: number }) {
		if (onBrowserClick) {
			onBrowserClick(item);
		}
	}

	function handleDeviceClick(item: { name: string; count: number }) {
		if (onDeviceClick) {
			onDeviceClick(item);
		}
	}

	function handleOsClick(item: { name: string; count: number }) {
		if (onOsClick) {
			onOsClick(item);
		}
	}
</script>

<div class="space-y-4">
	<!-- Tabs -->
	<div class="border-border flex gap-2 border-b">
		<Button
			variant={activeTab === 'browsers' ? 'default' : 'ghost'}
			size="sm"
			class="rounded-b-none"
			onclick={() => (activeTab = 'browsers')}
		>
			<Monitor class="mr-2 h-4 w-4" />
			Browsers
		</Button>
		<Button
			variant={activeTab === 'devices' ? 'default' : 'ghost'}
			size="sm"
			class="rounded-b-none"
			onclick={() => (activeTab = 'devices')}
		>
			<Smartphone class="mr-2 h-4 w-4" />
			Devices
		</Button>
		<Button
			variant={activeTab === 'os' ? 'default' : 'ghost'}
			size="sm"
			class="rounded-b-none"
			onclick={() => (activeTab = 'os')}
		>
			<Laptop class="mr-2 h-4 w-4" />
			OS
		</Button>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
			<div class="flex flex-col items-center gap-2">
				<div
					class="border-primary h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
				></div>
				<p class="text-sm">Loading data...</p>
			</div>
		</div>
	{:else if activeTab === 'browsers'}
		<TopItemsList
			items={browsers}
			labelKey="name"
			valueKey="count"
			maxItems={5}
			type="browser"
			showMoreTitle="All Browsers ({browsers.length} total)"
			onclick={handleBrowserClick}
		/>
	{:else if activeTab === 'devices'}
		<TopItemsList
			items={devices}
			labelKey="name"
			valueKey="count"
			maxItems={5}
			type="device"
			showMoreTitle="All Devices ({devices.length} total)"
			onclick={handleDeviceClick}
		/>
	{:else if activeTab === 'os'}
		<TopItemsList
			items={operatingSystems}
			labelKey="name"
			valueKey="count"
			maxItems={5}
			type="os"
			showMoreTitle="All Operating Systems ({operatingSystems.length} total)"
			onclick={handleOsClick}
		/>
	{/if}
</div>
