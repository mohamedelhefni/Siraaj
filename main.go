package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/mohamedelhefni/siraaj/geolocation"
)

//go:embed all:ui/dashboard
var uiFiles embed.FS

type Event struct {
	ID         uint64    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	EventName  string    `json:"event_name"`
	UserID     string    `json:"user_id"`
	SessionID  string    `json:"session_id"`
	URL        string    `json:"url"`
	Referrer   string    `json:"referrer"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	Country    string    `json:"country"`
	Browser    string    `json:"browser"`
	OS         string    `json:"os"`
	Device     string    `json:"device"`
	Properties string    `json:"properties"` // JSON string
	ProjectID  string    `json:"project_id"`
}

type Analytics struct {
	db *sql.DB
}

func NewAnalytics(dbPath string) (*Analytics, error) {
	// Use DuckDB with simplified settings to avoid WAL issues
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	// Enable DuckDB optimizations with error handling
	if _, err = db.Exec("PRAGMA memory_limit='512MB'"); err != nil {
		log.Printf("Warning: Could not set memory limit: %v", err)
	}

	if _, err = db.Exec("PRAGMA threads=2"); err != nil {
		log.Printf("Warning: Could not set threads: %v", err)
	}

	// Create tables with optimized schema
	// Use a simpler approach that's more compatible with DuckDB
	statements := []string{
		// Create table with auto-incrementing ID
		`CREATE TABLE IF NOT EXISTS events (
			id UBIGINT PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			event_name VARCHAR NOT NULL,
			user_id VARCHAR,
			session_id VARCHAR,
			url VARCHAR,
			referrer VARCHAR,
			user_agent VARCHAR,
			ip VARCHAR,
			country VARCHAR,
			browser VARCHAR,
			os VARCHAR,
			device VARCHAR,
			properties JSON,
			project_id VARCHAR DEFAULT 'default'
		)`,
		// Create sequence
		`CREATE SEQUENCE IF NOT EXISTS id_sequence START 1`,
		// Create indexes
		`CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_event_name ON events(event_name)`,
		`CREATE INDEX IF NOT EXISTS idx_user_id ON events(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_country ON events(country)`,
		`CREATE INDEX IF NOT EXISTS idx_referrer ON events(referrer)`,
		`CREATE INDEX IF NOT EXISTS idx_project_id ON events(project_id)`,
	}

	for i, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("Error executing statement %d: %v", i+1, err)
			log.Printf("Failed statement: %s", stmt)
			return nil, fmt.Errorf("failed to create database schema: %v", err)
		} else {
			log.Printf("✓ Successfully executed statement %d", i+1)
		}
	}

	return &Analytics{db: db}, nil
}

func (a *Analytics) Close() error {
	return a.db.Close()
}

func (a *Analytics) TrackEvent(event Event) error {
	query := `
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Default project_id if not provided
	if event.ProjectID == "" {
		event.ProjectID = "default"
	}

	_, err := a.db.Exec(query,
		event.Timestamp,
		event.EventName,
		event.UserID,
		event.SessionID,
		event.URL,
		event.Referrer,
		event.UserAgent,
		event.IP,
		event.Country,
		event.Browser,
		event.OS,
		event.Device,
		event.Properties,
		event.ProjectID,
	)

	if err != nil {
		log.Printf("Error inserting event: %v", err)
	}

	return err
}

// TrackEventBatch inserts multiple events efficiently
func (a *Analytics) TrackEventBatch(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		// Default project_id if not provided
		if event.ProjectID == "" {
			event.ProjectID = "default"
		}

		_, err := stmt.Exec(
			event.Timestamp,
			event.EventName,
			event.UserID,
			event.SessionID,
			event.URL,
			event.Referrer,
			event.UserAgent,
			event.IP,
			event.Country,
			event.Browser,
			event.OS,
			event.Device,
			event.Properties,
			event.ProjectID,
		)
		if err != nil {
			log.Printf("Error inserting batch event: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (a *Analytics) GetStats(startDate, endDate time.Time, limit int, filters map[string]string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	if limit <= 0 {
		limit = 10
	}

	// Build WHERE clause based on filters
	whereClause := "timestamp BETWEEN ? AND ?"
	args := []interface{}{startDate, endDate}

	if projectID, ok := filters["project"]; ok && projectID != "" {
		whereClause += " AND project_id = ?"
		args = append(args, projectID)
	}
	if source, ok := filters["source"]; ok && source != "" {
		whereClause += " AND referrer = ?"
		args = append(args, source)
	}
	if country, ok := filters["country"]; ok && country != "" {
		whereClause += " AND country = ?"
		args = append(args, country)
	}
	if browser, ok := filters["browser"]; ok && browser != "" {
		whereClause += " AND browser = ?"
		args = append(args, browser)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		whereClause += " AND event_name = ?"
		args = append(args, eventName)
	}

	// Use a single query with CTEs for better performance
	aggregateQuery := fmt.Sprintf(`
	WITH date_filtered AS (
		SELECT * FROM events 
		WHERE %s
	),
	event_stats AS (
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT session_id) as total_visits,
			COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_views,
			COUNT(DISTINCT CASE WHEN event_name = 'page_view' THEN session_id END) as sessions_with_views
		FROM date_filtered
	)
	SELECT total_events, unique_users, total_visits, page_views, sessions_with_views FROM event_stats;
	`, whereClause)

	var totalEvents, uniqueUsers, totalVisits, pageViews, sessionsWithViews int
	err := a.db.QueryRow(aggregateQuery, args...).Scan(&totalEvents, &uniqueUsers, &totalVisits, &pageViews, &sessionsWithViews)
	if err != nil {
		return nil, err
	}
	stats["total_events"] = totalEvents
	stats["unique_users"] = uniqueUsers
	stats["total_visits"] = totalVisits
	stats["page_views"] = pageViews

	// Calculate bounce rate: sessions with only 1 page view / total sessions
	var bounceRate float64
	if totalVisits > 0 {
		singlePageQuery := fmt.Sprintf(`
			SELECT COUNT(DISTINCT session_id) 
			FROM (
				SELECT session_id, COUNT(*) as view_count
				FROM events 
				WHERE %s AND event_name = 'page_view'
				GROUP BY session_id
				HAVING view_count = 1
			)
		`, whereClause)
		var singlePageSessions int
		err = a.db.QueryRow(singlePageQuery, args...).Scan(&singlePageSessions)
		if err == nil && sessionsWithViews > 0 {
			bounceRate = float64(singlePageSessions) / float64(sessionsWithViews) * 100
		}
	}
	stats["bounce_rate"] = bounceRate

	// Top events with optimized query
	query := fmt.Sprintf(`
		SELECT event_name, COUNT(*) as count 
		FROM events 
		WHERE %s
		GROUP BY event_name 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)
	queryArgs := append(args, limit)

	rows, err := a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topEvents := []map[string]interface{}{}
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			continue
		}
		topEvents = append(topEvents, map[string]interface{}{
			"name":  name,
			"count": count,
		})
	}
	stats["top_events"] = topEvents

	// Events over time using date_trunc for daily aggregation
	query = fmt.Sprintf(`
		SELECT 
			strftime(DATE_TRUNC('day', timestamp), '%%Y-%%m-%%d') as date, 
			COUNT(*) as count
		FROM events 
		WHERE %s
		GROUP BY date 
		ORDER BY date
	`, whereClause)

	rows, err = a.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timeline := []map[string]interface{}{}
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			continue
		}
		timeline = append(timeline, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	stats["timeline"] = timeline

	// Top pages
	query = fmt.Sprintf(`
		SELECT url, COUNT(*) as count 
		FROM events 
		WHERE %s AND url IS NOT NULL AND url != ''
		GROUP BY url 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topPages := []map[string]interface{}{}
	for rows.Next() {
		var url string
		var count int
		if err := rows.Scan(&url, &count); err != nil {
			continue
		}
		topPages = append(topPages, map[string]interface{}{
			"url":   url,
			"count": count,
		})
	}
	stats["top_pages"] = topPages

	// Browsers
	query = fmt.Sprintf(`
		SELECT browser, COUNT(*) as count 
		FROM events 
		WHERE %s AND browser IS NOT NULL AND browser != ''
		GROUP BY browser 
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	rows, err = a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	browsers := []map[string]interface{}{}
	for rows.Next() {
		var browser string
		var count int
		if err := rows.Scan(&browser, &count); err != nil {
			continue
		}
		browsers = append(browsers, map[string]interface{}{
			"name":  browser,
			"count": count,
		})
	}
	stats["browsers"] = browsers

	// Top Countries
	query = fmt.Sprintf(`
		SELECT country, COUNT(*) as count 
		FROM events 
		WHERE %s AND country IS NOT NULL AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topCountries := []map[string]interface{}{}
	for rows.Next() {
		var country string
		var count int
		if err := rows.Scan(&country, &count); err != nil {
			continue
		}
		topCountries = append(topCountries, map[string]interface{}{
			"name":  country,
			"count": count,
		})
	}
	stats["top_countries"] = topCountries

	// Top Sources (Referrers) with URL parsing
	query = fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN referrer = '' OR referrer IS NULL THEN 'Direct'
				ELSE referrer
			END as source,
			COUNT(*) as count 
		FROM events 
		WHERE %s
		GROUP BY source 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = a.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topSources := []map[string]interface{}{}
	for rows.Next() {
		var referrer string
		var count int
		if err := rows.Scan(&referrer, &count); err != nil {
			continue
		}
		topSources = append(topSources, map[string]interface{}{
			"name":  referrer,
			"count": count,
		})
	}
	stats["top_sources"] = topSources

	// Calculate trends by comparing with previous period
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := startDate

	prevWhereClause := "timestamp BETWEEN ? AND ?"
	prevArgs := []interface{}{prevStartDate, prevEndDate}

	// Apply same filters to previous period
	if projectID, ok := filters["project"]; ok && projectID != "" {
		prevWhereClause += " AND project_id = ?"
		prevArgs = append(prevArgs, projectID)
	}
	if source, ok := filters["source"]; ok && source != "" {
		prevWhereClause += " AND referrer = ?"
		prevArgs = append(prevArgs, source)
	}
	if country, ok := filters["country"]; ok && country != "" {
		prevWhereClause += " AND country = ?"
		prevArgs = append(prevArgs, country)
	}
	if browser, ok := filters["browser"]; ok && browser != "" {
		prevWhereClause += " AND browser = ?"
		prevArgs = append(prevArgs, browser)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		prevWhereClause += " AND event_name = ?"
		prevArgs = append(prevArgs, eventName)
	}

	prevQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT session_id) as total_visits,
			COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_views
		FROM events 
		WHERE %s
	`, prevWhereClause)

	var prevTotalEvents, prevUniqueUsers, prevTotalVisits, prevPageViews int
	err = a.db.QueryRow(prevQuery, prevArgs...).Scan(&prevTotalEvents, &prevUniqueUsers, &prevTotalVisits, &prevPageViews)
	if err == nil {
		stats["prev_total_events"] = prevTotalEvents
		stats["prev_unique_users"] = prevUniqueUsers
		stats["prev_total_visits"] = prevTotalVisits
		stats["prev_page_views"] = prevPageViews

		// Calculate percentage changes
		if prevTotalEvents > 0 {
			stats["events_change"] = float64(totalEvents-prevTotalEvents) / float64(prevTotalEvents) * 100
		}
		if prevUniqueUsers > 0 {
			stats["users_change"] = float64(uniqueUsers-prevUniqueUsers) / float64(prevUniqueUsers) * 100
		}
		if prevTotalVisits > 0 {
			stats["visits_change"] = float64(totalVisits-prevTotalVisits) / float64(prevTotalVisits) * 100
		}
		if prevPageViews > 0 {
			stats["page_views_change"] = float64(pageViews-prevPageViews) / float64(prevPageViews) * 100
		}
	}

	return stats, nil
}

// GetEvents returns paginated events
func (a *Analytics) GetEvents(startDate, endDate time.Time, limit, offset int) (map[string]interface{}, error) {
	query := `
		SELECT id, timestamp, event_name, user_id, session_id, url, referrer, 
			   user_agent, ip, country, browser, os, device, properties
		FROM events 
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := a.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Timestamp,
			&event.EventName,
			&event.UserID,
			&event.SessionID,
			&event.URL,
			&event.Referrer,
			&event.UserAgent,
			&event.IP,
			&event.Country,
			&event.Browser,
			&event.OS,
			&event.Device,
			&event.Properties,
		)
		if err != nil {
			continue
		}
		events = append(events, event)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM events WHERE timestamp BETWEEN ? AND ?`
	err = a.db.QueryRow(countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, nil
}

// GetOnlineUsers returns currently online users (active in last N minutes)
func (a *Analytics) GetOnlineUsers(timeWindowMinutes int) (map[string]interface{}, error) {
	cutoffTime := time.Now().Add(-time.Duration(timeWindowMinutes) * time.Minute)

	query := `
		SELECT 
			COUNT(DISTINCT user_id) as online_users,
			COUNT(DISTINCT session_id) as active_sessions
		FROM events 
		WHERE timestamp >= ?
	`

	var onlineUsers, activeSessions int
	err := a.db.QueryRow(query, cutoffTime).Scan(&onlineUsers, &activeSessions)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"online_users":     onlineUsers,
		"active_sessions":  activeSessions,
		"time_window_mins": timeWindowMinutes,
		"cutoff_time":      cutoffTime,
	}, nil
}

// GetProjects returns list of distinct project IDs
func (a *Analytics) GetProjects() ([]string, error) {
	query := `SELECT DISTINCT project_id FROM events WHERE project_id IS NOT NULL AND project_id != '' ORDER BY project_id`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []string
	for rows.Next() {
		var projectID string
		if err := rows.Scan(&projectID); err != nil {
			continue
		}
		projects = append(projects, projectID)
	}

	return projects, nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set allowed origin(s)
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace with your frontend origin
		// Set allowed HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Set allowed request headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Allow credentials (if using cookies or HTTP authentication)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests (OPTIONS method)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ResponseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log the request
		duration := time.Since(start)
		log.Printf("[%s] %s %s - Status: %d - Duration: %v - IP: %s - User-Agent: %s",
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)
	})
}

// getClientIP extracts the real client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func main() {
	// Initialize geolocation service
	geoService, err := geolocation.NewService()
	if err != nil {
		log.Printf("⚠️  Warning: Geolocation service unavailable: %v", err)
		log.Println("⚠️  Continuing without geolocation support...")
		geoService = nil
	}
	if geoService != nil {
		defer geoService.Close()
	}

	analytics, err := NewAnalytics("analytics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer analytics.Close()

	log.Println("✓ DuckDB initialized successfully")

	// Track event endpoint
	http.HandleFunc("/api/track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var event Event
		if err := json.Unmarshal(body, &event); err != nil {
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
		if geoService != nil && event.Country == "" {
			if geo := geoService.LookupOrDefault(event.IP); geo != nil {
				event.Country = geo.Country
				if event.Country == "" {
					event.Country = geo.CountryCode
				}
			}
		}

		if err := analytics.TrackEvent(event); err != nil {
			log.Printf("Error tracking event: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Stats endpoint
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
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

		// First, let's check if there are any events at all
		var totalCount int
		err := analytics.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalCount)
		if err != nil {
			log.Printf("❌ Error counting total events: %v", err)
		} else {
			log.Printf("📊 Total events in database: %d", totalCount)
		}

		// Check events in date range
		var rangeCount int
		err = analytics.db.QueryRow("SELECT COUNT(*) FROM events WHERE timestamp BETWEEN ? AND ?", startDate, endDate).Scan(&rangeCount)
		if err != nil {
			log.Printf("❌ Error counting events in range: %v", err)
		} else {
			log.Printf("📊 Events in date range (%s to %s): %d",
				startDate.Format("2006-01-02"),
				endDate.Format("2006-01-02"),
				rangeCount)
		}

		// Debug: Show min and max timestamps in DB
		var minTimestamp, maxTimestamp time.Time
		err = analytics.db.QueryRow("SELECT MIN(timestamp), MAX(timestamp) FROM events").Scan(&minTimestamp, &maxTimestamp)
		if err == nil {
			log.Printf("📊 Event date range in DB: %s to %s",
				minTimestamp.Format("2006-01-02 15:04:05"),
				maxTimestamp.Format("2006-01-02 15:04:05"))
		}

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

		stats, err := analytics.GetStats(startDate, endDate, limit, filters)
		if err != nil {
			log.Printf("Error getting stats: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Events endpoint with pagination
	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
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

		events, err := analytics.GetEvents(startDate, endDate, limit, offset)
		if err != nil {
			log.Printf("Error getting events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})

	// Online users endpoint
	http.HandleFunc("/api/online", func(w http.ResponseWriter, r *http.Request) {
		timeWindow := 5 // default 5 minutes
		if windowStr := r.URL.Query().Get("window"); windowStr != "" {
			fmt.Sscanf(windowStr, "%d", &timeWindow)
			if timeWindow > 60 {
				timeWindow = 60 // cap at 1 hour
			}
		}

		online, err := analytics.GetOnlineUsers(timeWindow)
		if err != nil {
			log.Printf("Error getting online users: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(online)
	})

	// Projects endpoint - list all projects
	http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		projects, err := analytics.GetProjects()
		if err != nil {
			log.Printf("Error getting projects: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	})

	// Debug endpoint to show all events
	http.HandleFunc("/api/debug/events", func(w http.ResponseWriter, r *http.Request) {
		rows, err := analytics.db.Query("SELECT id, timestamp, event_name, user_id FROM events ORDER BY timestamp DESC LIMIT 50")
		if err != nil {
			log.Printf("Error querying events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := []map[string]interface{}{}
		for rows.Next() {
			var id uint64
			var timestamp time.Time
			var eventName, userID string
			if err := rows.Scan(&id, &timestamp, &eventName, &userID); err != nil {
				continue
			}
			events = append(events, map[string]interface{}{
				"id":         id,
				"timestamp":  timestamp.Format(time.RFC3339),
				"event_name": eventName,
				"user_id":    userID,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"events": events,
			"count":  len(events),
		})
	})

	// Health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "ok",
			"database":    "duckdb",
			"version":     "1.0.0",
			"geolocation": geoService != nil,
		})
	})

	// Geolocation test endpoint
	http.HandleFunc("/api/geo", func(w http.ResponseWriter, r *http.Request) {
		if geoService == nil {
			http.Error(w, "Geolocation service not available", http.StatusServiceUnavailable)
			return
		}

		ip := r.URL.Query().Get("ip")
		if ip == "" {
			ip = getClientIP(r)
		}

		geo := geoService.LookupOrDefault(ip)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ip":           ip,
			"country":      geo.Country,
			"country_code": geo.CountryCode,
			"city":         geo.City,
		})
	})

	// Serve dashboard (SvelteKit app)
	dashboardFS, err := fs.Sub(uiFiles, "ui/dashboard")
	if err != nil {
		log.Printf("Warning: Could not load dashboard: %v", err)
	} else {
		http.Handle("/dashboard/", http.StripPrefix("/dashboard", http.FileServer(http.FS(dashboardFS))))
	}

	// Serve UI (must be last as it's a catch-all)
	http.Handle("/", http.FileServer(http.FS(uiFiles)))

	port := "8080"
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Analytics Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🎨 Dashboard:  http://localhost:%s/dashboard/\n", port)
	fmt.Printf("📡 API Track:  http://localhost:%s/api/track\n", port)
	fmt.Printf("📈 API Stats:  http://localhost:%s/api/stats\n", port)
	fmt.Printf("🌍 Geo Test:   http://localhost:%s/api/geo\n", port)
	fmt.Printf("❤️  Health:    http://localhost:%s/api/health\n", port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Server ready - Using official DuckDB Go driver")
	fmt.Println("✓ Svelte Dashboard embedded and ready")
	if geoService != nil {
		fmt.Println("✓ Geolocation service enabled")
	} else {
		fmt.Println("⚠️  Geolocation service disabled")
	}
	fmt.Println()

	// Apply middleware: logging first, then CORS
	handler := loggingMiddleware(enableCORS(http.DefaultServeMux))
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
