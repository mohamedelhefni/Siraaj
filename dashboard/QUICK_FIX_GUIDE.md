# Quick Fix Guide: Dashboard Loading States

Run these find-and-replace operations in VS Code on `dashboard/src/routes/+page.svelte`:

## 1. Pages Section (line ~1030)
**Find exactly:**
```
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
```

**Replace with:**
```
{#if pagesLoading || entryExitLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
						<div class="flex flex-col items-center gap-2">
							<div
								class="border-primary h-6 w-6 animate-spin rounded-full border-2 border-t-transparent"
							></div>
							<p class="text-xs">Loading...</p>
						</div>
					</div>
				{:else if pagesTab === 'all'}
```

## 2. Countries Panel (line ~1078)  
**Find:** `loading={statsLoading}` (in CountriesPanel component)
**Replace:** `loading={countriesLoading}`

## 3. Sources Section (line ~1090)
**Find exactly:**
```
				{#if statsLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
```

**Replace with:**
```
				{#if sourcesLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
```

## 4. Events Section (line ~1120)
**Find exactly:**
```
				{#if statsLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
```

**Replace with:**
```
				{#if eventsLoading}
					<div class="text-muted-foreground flex min-h-[200px] items-center justify-center">
```

## 5. Devices/Browsers Panel (line ~1155)
**Find:** `loading={statsLoading}` (in BrowserPanel or similar component near line 1155)
**Replace:** `loading={devicesLoading}`

---

## Verification
After all replacements, search for `statsLoading` in the file.
You should find:
- ✅ Line ~102: Declaration `let topStatsLoading = $state(false);`
- ✅ Line ~327: Assignment `topStatsLoading = true;`  
- ✅ Line ~350: Assignment `topStatsLoading = false;`
- ❌ NO other references

Total: 3 references only (all for the new `topStatsLoading` variable).

## Test
```bash
cd dashboard
pnpm run dev
```

Open the dashboard and watch the Network tab - you should see:
1. Online users card appears almost instantly (~50ms)
2. Top pages/countries appear quickly (~200-500ms)
3. Timeline and slow queries appear last (~800-2000ms)
4. Each section shows its own loading spinner independently
