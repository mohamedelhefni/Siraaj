package service

import (
	"log"
	"time"

	"github.com/mohamedelhefni/siraaj/internal/repository"
)

type CacheWarmer struct {
	repo repository.EventRepository
}

func NewCacheWarmer(repo repository.EventRepository) *CacheWarmer {
	return &CacheWarmer{repo: repo}
}

// Start begins the cache warming process in the background
func (w *CacheWarmer) Start() {
	// Warm cache every 5 minutes for common queries
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Initial warm on startup
	log.Println("[CacheWarmer] Starting initial cache warm...")
	w.warmCommonQueries()

	for range ticker.C {
		w.warmCommonQueries()
	}
}

// warmCommonQueries pre-loads cache with frequently accessed queries
func (w *CacheWarmer) warmCommonQueries() {
	log.Println("[CacheWarmer] Warming cache for common queries...")

	now := time.Now()

	// Common date ranges that users typically view
	dateRanges := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{"today", time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), now},
		{"yesterday", now.AddDate(0, 0, -1).Truncate(24 * time.Hour), now.AddDate(0, 0, -1).Add(23*time.Hour + 59*time.Minute)},
		{"last_7_days", now.AddDate(0, 0, -7), now},
		{"last_30_days", now.AddDate(0, 0, -30), now},
		{"this_month", time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), now},
	}

	// Get all projects to warm cache for each
	projects, err := w.repo.GetProjects()
	if err != nil {
		log.Printf("[CacheWarmer] Error getting projects: %v", err)
		return
	}

	warmCount := 0

	for _, dr := range dateRanges {
		// Warm cache for overall stats (no project filter)
		filters := make(map[string]string)
		w.warmForFilters(dr.start, dr.end, filters)
		warmCount++

		// Warm cache for each project
		for _, project := range projects {
			projectFilters := map[string]string{"project": project}
			w.warmForFilters(dr.start, dr.end, projectFilters)
			warmCount++
		}
	}

	log.Printf("[CacheWarmer] Cache warming completed (%d queries warmed)", warmCount)
}

// warmForFilters warms the cache for a specific set of filters
func (w *CacheWarmer) warmForFilters(startDate, endDate time.Time, filters map[string]string) {
	// Warm top stats (most frequently accessed)
	_, err := w.repo.GetTopStats(startDate, endDate, filters)
	if err != nil {
		log.Printf("[CacheWarmer] Error warming GetTopStats: %v", err)
	}

	// Warm timeline (second most frequent)
	_, err = w.repo.GetTimeline(startDate, endDate, filters)
	if err != nil {
		log.Printf("[CacheWarmer] Error warming GetTimeline: %v", err)
	}

	// Warm top pages
	_, err = w.repo.GetTopPages(startDate, endDate, 10, filters)
	if err != nil {
		log.Printf("[CacheWarmer] Error warming GetTopPages: %v", err)
	}

	// Warm top countries
	_, err = w.repo.GetTopCountries(startDate, endDate, 10, filters)
	if err != nil {
		log.Printf("[CacheWarmer] Error warming GetTopCountries: %v", err)
	}

	// Warm browsers/devices/os
	_, err = w.repo.GetBrowsersDevicesOS(startDate, endDate, 5, filters)
	if err != nil {
		log.Printf("[CacheWarmer] Error warming GetBrowsersDevicesOS: %v", err)
	}
}
