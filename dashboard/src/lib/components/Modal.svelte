<script>
	import { X } from 'lucide-svelte';

	let { open = $bindable(false), title = '', children } = $props();

	function handleClose() {
		open = false;
	}

	function handleBackdropClick(e) {
		if (e.target === e.currentTarget) {
			handleClose();
		}
	}

	function handleKeydown(e) {
		if (e.key === 'Escape') {
			handleClose();
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
		tabindex="-1"
	>
		<div
			class="bg-background relative max-h-[90vh] w-full max-w-2xl overflow-hidden rounded-lg shadow-lg"
			role="document"
		>
			<!-- Header -->
			<div class="border-b px-6 py-4">
				<div class="flex items-center justify-between">
					<h2 id="modal-title" class="text-xl font-semibold">
						{title}
					</h2>
					<button
						type="button"
						onclick={handleClose}
						class="hover:bg-accent rounded-md p-1 transition-colors"
						aria-label="Close"
					>
						<X class="h-5 w-5" />
					</button>
				</div>
			</div>

			<!-- Content -->
			<div class="max-h-[calc(90vh-8rem)] overflow-y-auto px-6 py-4">
				{@render children()}
			</div>

			<!-- Footer (optional) -->
			<div class="border-t px-6 py-4">
				<button
					type="button"
					onclick={handleClose}
					class="bg-primary text-primary-foreground hover:bg-primary/90 rounded-md px-4 py-2 text-sm font-medium transition-colors"
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}
