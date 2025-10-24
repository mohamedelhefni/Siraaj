package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

//go:embed ui/*
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
}

type Analytics struct {
	db *sql.DB
}

func NewAnalytics(dbPath string) (*Analytics, error) {
	// Use DuckDB with optimized settings
	connStr := dbPath + "?access_mode=read_write&threads=4"
	db, err := sql.Open("duckdb", connStr)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Enable DuckDB optimizations
	_, err = db.Exec("PRAGMA memory_limit='1GB'")
	if err != nil {
		log.Printf("Warning: Could not set memory limit: %v", err)
	}

	_, err = db.Exec("PRAGMA threads=4")
	if err != nil {
		log.Printf("Warning: Could not set threads: %v", err)
	}

	// Create tables with optimized schema
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id UBIGINT PRIMARY KEY ,
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
		properties JSON
	);

	CREATE SEQUENCE id_sequence START 1;
	ALTER TABLE events ALTER COLUMN id SET DEFAULT NEXTVAL('id_sequence');


	CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_event_name ON events(event_name);
	CREATE INDEX IF NOT EXISTS idx_user_id ON events(user_id);
	CREATE INDEX IF NOT EXISTS idx_country ON events(country);
	CREATE INDEX IF NOT EXISTS idx_referrer ON events(referrer);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &Analytics{db: db}, nil
}

func (a *Analytics) Close() error {
	return a.db.Close()
}

func (a *Analytics) TrackEvent(event Event) error {
	query := `
		INSERT INTO events (timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

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
	)

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
		INSERT INTO events (timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
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
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (a *Analytics) GetStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Use a single query with CTEs for better performance
	aggregateQuery := `
	WITH date_filtered AS (
		SELECT * FROM events 
		WHERE timestamp BETWEEN ? AND ?
	),
	event_stats AS (
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users
		FROM date_filtered
	)
	SELECT total_events, unique_users FROM event_stats;
	`

	var totalEvents, uniqueUsers int
	err := a.db.QueryRow(aggregateQuery, startDate, endDate).Scan(&totalEvents, &uniqueUsers)
	if err != nil {
		return nil, err
	}
	stats["total_events"] = totalEvents
	stats["unique_users"] = uniqueUsers

	// Top events with optimized query
	rows, err := a.db.Query(`
		SELECT event_name, COUNT(*) as count 
		FROM events 
		WHERE timestamp BETWEEN ? AND ?
		GROUP BY event_name 
		ORDER BY count DESC 
		LIMIT 10
	`, startDate, endDate)
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
	rows, err = a.db.Query(`
		SELECT 
			strftime(DATE_TRUNC('day', timestamp), '%Y-%m-%d') as date, 
			COUNT(*) as count
		FROM events 
		WHERE timestamp BETWEEN ? AND ?
		GROUP BY date 
		ORDER BY date
	`, startDate, endDate)
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
	rows, err = a.db.Query(`
		SELECT url, COUNT(*) as count 
		FROM events 
		WHERE timestamp BETWEEN ? AND ? AND url IS NOT NULL AND url != ''
		GROUP BY url 
		ORDER BY count DESC 
		LIMIT 10
	`, startDate, endDate)
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
	rows, err = a.db.Query(`
		SELECT browser, COUNT(*) as count 
		FROM events 
		WHERE timestamp BETWEEN ? AND ? AND browser IS NOT NULL AND browser != ''
		GROUP BY browser 
		ORDER BY count DESC
	`, startDate, endDate)
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
	rows, err = a.db.Query(`
		SELECT country, COUNT(*) as count 
		FROM events 
		WHERE timestamp BETWEEN ? AND ? AND country IS NOT NULL AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT 10
	`, startDate, endDate)
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
	rows, err = a.db.Query(`
		SELECT 
			CASE 
				WHEN referrer = '' OR referrer IS NULL THEN 'Direct'
				ELSE referrer
			END as source,
			COUNT(*) as count 
		FROM events 
		WHERE timestamp BETWEEN ? AND ?
		GROUP BY source 
		ORDER BY count DESC 
		LIMIT 10
	`, startDate, endDate)
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

	return stats, nil
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

func main() {
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
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var event Event
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		// Get IP from request
		if event.IP == "" {
			event.IP = r.RemoteAddr
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
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7)

		// Parse date range from query params
		if start := r.URL.Query().Get("start"); start != "" {
			if t, err := time.Parse("2006-01-02", start); err == nil {
				startDate = t
			}
		}
		if end := r.URL.Query().Get("end"); end != "" {
			if t, err := time.Parse("2006-01-02", end); err == nil {
				endDate = t
			}
		}

		stats, err := analytics.GetStats(startDate, endDate)
		if err != nil {
			log.Printf("Error getting stats: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Serve UI
	http.Handle("/", http.FileServer(http.FS(uiFiles)))

	// Health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "ok",
			"database": "duckdb",
			"version":  "1.0.0",
		})
	})

	port := "8080"
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Analytics Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🌐 Dashboard:  http://localhost:%s\n", port)
	fmt.Printf("📡 API Track:  http://localhost:%s/api/track\n", port)
	fmt.Printf("📈 API Stats:  http://localhost:%s/api/stats\n", port)
	fmt.Printf("❤️  Health:     http://localhost:%s/api/health\n", port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Server ready - Using official DuckDB Go driver")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux)))
}
