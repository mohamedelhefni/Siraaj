package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedelhefni/siraaj/geolocation"
	"github.com/mohamedelhefni/siraaj/internal/domain"
	"github.com/mohamedelhefni/siraaj/internal/service"
)

type EventHandler struct {
	service    service.EventService
	geoService *geolocation.Service
}

func NewEventHandler(service service.EventService, geoService *geolocation.Service) *EventHandler {
	return &EventHandler{
		service:    service,
		geoService: geoService,
	}
}

func (h *EventHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("Error Unmarshal json: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Get IP from request
	if event.IP == "" {
		event.IP = getClientIP(r)
	}

	// Enrich with geolocation data if service is available
	if h.geoService != nil && event.Country == "" {
		if geo := h.geoService.LookupOrDefault(event.IP); geo != nil {
			event.Country = geo.Country
			if event.Country == "" {
				event.Country = geo.CountryCode
			}
		}
	}

	if err := h.service.TrackEvent(event); err != nil {
		log.Printf("Error tracking event: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *EventHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// Default to last 7 days
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, 0, -7)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Parse date range from query params
	if start := r.URL.Query().Get("start"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			// Set to beginning of day for start date
			startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			// Set to end of day for the end date
			endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		}
	}

	log.Printf("📅 Stats query: startDate=%v, endDate=%v", startDate, endDate)
	log.Printf("📅 Date range: %s to %s", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

	// Parse limit parameter
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l == 1 {
			if limit > 1000 {
				limit = 1000 // Cap at 1000
			}
		}
	}

	// Parse filters
	filters := make(map[string]string)
	if project := r.URL.Query().Get("project"); project != "" {
		filters["project"] = project
	}
	if source := r.URL.Query().Get("source"); source != "" {
		filters["source"] = source
	}
	if country := r.URL.Query().Get("country"); country != "" {
		filters["country"] = country
	}
	if browser := r.URL.Query().Get("browser"); browser != "" {
		filters["browser"] = browser
	}
	if event := r.URL.Query().Get("event"); event != "" {
		filters["event"] = event
	}
	if metric := r.URL.Query().Get("metric"); metric != "" {
		filters["metric"] = metric
	}

	stats, err := h.service.GetStats(startDate, endDate, limit, filters)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse date range
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, 0, -7)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	if start := r.URL.Query().Get("start"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		}
	}

	// Parse pagination parameters
	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l == 1 {
			if limit > 1000 {
				limit = 1000
			}
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}

	events, err := h.service.GetEvents(startDate, endDate, limit, offset)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	timeWindow := 5
	if windowStr := r.URL.Query().Get("window"); windowStr != "" {
		fmt.Sscanf(windowStr, "%d", &timeWindow)
		if timeWindow > 60 {
			timeWindow = 60
		}
	}

	online, err := h.service.GetOnlineUsers(timeWindow)
	if err != nil {
		log.Printf("Error getting online users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(online)
}

func (h *EventHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetProjects()
	if err != nil {
		log.Printf("Error getting projects: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *EventHandler) GetTopProperties(w http.ResponseWriter, r *http.Request) {
	// Default to last 7 days
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, 0, -7)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Parse date range from query params
	if start := r.URL.Query().Get("start"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
	}
	if end := r.URL.Query().Get("end"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		}
	}

	// Parse limit parameter
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l == 1 {
			if limit > 100 {
				limit = 100 // Cap at 100
			}
		}
	}

	// Parse filters
	filters := make(map[string]string)
	if project := r.URL.Query().Get("project"); project != "" {
		filters["project"] = project
	}
	if source := r.URL.Query().Get("source"); source != "" {
		filters["source"] = source
	}
	if country := r.URL.Query().Get("country"); country != "" {
		filters["country"] = country
	}
	if browser := r.URL.Query().Get("browser"); browser != "" {
		filters["browser"] = browser
	}
	if event := r.URL.Query().Get("event"); event != "" {
		filters["event"] = event
	}

	properties, err := h.service.GetTopProperties(startDate, endDate, limit, filters)
	if err != nil {
		log.Printf("Error getting top properties: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}

func (h *EventHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"database":    "duckdb",
		"version":     "1.0.0",
		"geolocation": h.geoService != nil,
	})
}

func (h *EventHandler) GeoTest(w http.ResponseWriter, r *http.Request) {
	if h.geoService == nil {
		http.Error(w, "Geolocation service not available", http.StatusServiceUnavailable)
		return
	}

	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = getClientIP(r)
	}

	geo := h.geoService.LookupOrDefault(ip)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ip":           ip,
		"country":      geo.Country,
		"country_code": geo.CountryCode,
		"city":         geo.City,
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}

func parseBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	browsers := map[string]string{
		"edg":     "Edge",
		"chrome":  "Chrome",
		"safari":  "Safari",
		"firefox": "Firefox",
		"opera":   "Opera",
	}

	for key, name := range browsers {
		if strings.Contains(ua, key) {
			if key == "safari" && strings.Contains(ua, "chrome") {
				continue
			}
			return name
		}
	}
	return "Other"
}

func parseOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	osList := map[string]string{
		"windows": "Windows",
		"mac":     "macOS",
		"linux":   "Linux",
		"android": "Android",
		"iphone":  "iOS",
		"ipad":    "iOS",
	}

	for key, name := range osList {
		if strings.Contains(ua, key) {
			return name
		}
	}
	return "Other"
}

func parseDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") {
		return "Mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "Tablet"
	}
	return "Desktop"
}
