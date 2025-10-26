package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mohamedelhefni/siraaj/geolocation"
	"github.com/mohamedelhefni/siraaj/internal/botdetector"
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
		geo := h.geoService.LookupOrDefault(event.IP)
		if geo != nil {
			event.Country = geo.Country
			if event.Country == "" {
				event.Country = geo.CountryCode
			}
		}
	}

	// Detect if user agent belongs to a bot
	event.IsBot = botdetector.IsBot(event.UserAgent)
	if event.IsBot {
		log.Printf("🤖 Bot detected: %s", botdetector.GetBotName(event.UserAgent))
	}

	if err := h.service.TrackEvent(event); err != nil {
		log.Printf("Error tracking event: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
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
		var l int
		if n, err := fmt.Sscanf(limitStr, "%d", &l); err == nil && n == 1 {
			limit = l
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
	if botFilter := r.URL.Query().Get("botFilter"); botFilter != "" {
		filters["botFilter"] = botFilter
	}

	stats, err := h.service.GetStats(startDate, endDate, limit, filters)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats: %v", err)
	}
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
		var l int
		if n, err := fmt.Sscanf(limitStr, "%d", &l); err == nil && n == 1 {
			limit = l
			if limit > 1000 {
				limit = 1000
			}
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		var o int
		if _, err := fmt.Sscanf(offsetStr, "%d", &o); err == nil {
			offset = o
		}
	}

	events, err := h.service.GetEvents(startDate, endDate, limit, offset)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("Error encoding events: %v", err)
	}
}

func (h *EventHandler) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	timeWindow := 5
	if windowStr := r.URL.Query().Get("window"); windowStr != "" {
		var tw int
		if _, err := fmt.Sscanf(windowStr, "%d", &tw); err == nil {
			timeWindow = tw
			if timeWindow > 60 {
				timeWindow = 60
			}
		}
	}

	online, err := h.service.GetOnlineUsers(timeWindow)
	if err != nil {
		log.Printf("Error getting online users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(online); err != nil {
		log.Printf("Error encoding online users: %v", err)
	}
}

func (h *EventHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetProjects()
	if err != nil {
		log.Printf("Error getting projects: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		log.Printf("Error encoding projects: %v", err)
	}
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
		var l int
		if n, err := fmt.Sscanf(limitStr, "%d", &l); err == nil && n == 1 {
			limit = l
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
	if err := json.NewEncoder(w).Encode(properties); err != nil {
		log.Printf("Error encoding properties: %v", err)
	}
}

func (h *EventHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"database":    "duckdb",
		"version":     "1.0.0",
		"geolocation": h.geoService != nil,
	}); err != nil {
		log.Printf("Error encoding health response: %v", err)
	}
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ip":           ip,
		"country":      geo.Country,
		"country_code": geo.CountryCode,
		"city":         geo.City,
	}); err != nil {
		log.Printf("Error encoding geo response: %v", err)
	}
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
