package main

// Funnel Data Generator for Load Testing
//
// This script generates realistic user journey data for testing funnel analysis.
// It creates sequential events that form common funnel patterns like:
// - E-commerce: Homepage → Product → Cart → Checkout → Purchase
// - SaaS: Landing → Signup → Onboarding → Feature Use → Upgrade
// - Content: Homepage → Article → Newsletter → Download
//
// Usage:
//   go run funnel_data.go <mode> <num_users> [db_path]
//
// Modes:
//   db   - Insert directly into DuckDB (fastest)
//   http - Send via HTTP API (realistic)
//
// Examples:
//   go run funnel_data.go db 10000              # Generate 10k user journeys
//   go run funnel_data.go http 5000             # 5k journeys via HTTP
//   go run funnel_data.go db 50000 ../data/analytics.db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

// Event represents an analytics event
type Event struct {
	Timestamp       time.Time `json:"timestamp"`
	EventName       string    `json:"event_name"`
	UserID          string    `json:"user_id"`
	SessionID       string    `json:"session_id"`
	SessionDuration int       `json:"session_duration"`
	URL             string    `json:"url"`
	Referrer        string    `json:"referrer"`
	UserAgent       string    `json:"user_agent"`
	IP              string    `json:"ip"`
	Country         string    `json:"country"`
	Browser         string    `json:"browser"`
	OS              string    `json:"os"`
	Device          string    `json:"device"`
	IsBot           bool      `json:"is_bot"`
	ProjectID       string    `json:"project_id"`
}

// FunnelStep represents a step in a user journey
type FunnelStep struct {
	EventName   string
	URL         string
	MinDuration int // Minimum seconds to next step
	MaxDuration int // Maximum seconds to next step
	DropOffRate float64
}

// FunnelTemplate defines a complete user journey
type FunnelTemplate struct {
	Name    string
	Steps   []FunnelStep
	Weight  int // Higher = more likely to be chosen
	Country string
	Browser string
	Device  string
}

// Funnel templates for realistic user journeys
var funnelTemplates = []FunnelTemplate{
	// E-commerce: Complete Purchase Journey
	{
		Name:    "E-commerce Purchase",
		Weight:  30,
		Country: "United States",
		Browser: "Chrome",
		Device:  "Desktop",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.20},
			{EventName: "page_view", URL: "/product/123", MinDuration: 10, MaxDuration: 60, DropOffRate: 0.30},
			{EventName: "add_to_cart", URL: "/product/123", MinDuration: 2, MaxDuration: 10, DropOffRate: 0.40},
			{EventName: "page_view", URL: "/cart", MinDuration: 5, MaxDuration: 20, DropOffRate: 0.35},
			{EventName: "checkout_started", URL: "/checkout", MinDuration: 10, MaxDuration: 120, DropOffRate: 0.50},
			{EventName: "payment_completed", URL: "/payment", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.10},
			{EventName: "purchase", URL: "/confirmation", MinDuration: 3, MaxDuration: 10, DropOffRate: 0.05},
		},
	},
	// E-commerce: Browse Only (High Drop-off)
	{
		Name:    "E-commerce Browse",
		Weight:  40,
		Country: "United States",
		Browser: "Safari",
		Device:  "Mobile",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/", MinDuration: 3, MaxDuration: 15, DropOffRate: 0.30},
			{EventName: "page_view", URL: "/product/456", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.70},
			{EventName: "page_view", URL: "/product/789", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.80},
		},
	},
	// SaaS: Signup to Activation
	{
		Name:    "SaaS Activation",
		Weight:  25,
		Country: "United Kingdom",
		Browser: "Chrome",
		Device:  "Desktop",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/", MinDuration: 10, MaxDuration: 60, DropOffRate: 0.25},
			{EventName: "page_view", URL: "/pricing", MinDuration: 15, MaxDuration: 90, DropOffRate: 0.40},
			{EventName: "button_click", URL: "/pricing", MinDuration: 2, MaxDuration: 5, DropOffRate: 0.15},
			{EventName: "signup", URL: "/signup", MinDuration: 30, MaxDuration: 180, DropOffRate: 0.45},
			{EventName: "page_view", URL: "/dashboard", MinDuration: 5, MaxDuration: 20, DropOffRate: 0.30},
			{EventName: "feature_used", URL: "/dashboard", MinDuration: 10, MaxDuration: 300, DropOffRate: 0.25},
		},
	},
	// Content: Newsletter Signup
	{
		Name:    "Content Newsletter",
		Weight:  35,
		Country: "Germany",
		Browser: "Firefox",
		Device:  "Desktop",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/blog", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.15},
			{EventName: "page_view", URL: "/blog/article-1", MinDuration: 60, MaxDuration: 300, DropOffRate: 0.40},
			{EventName: "button_click", URL: "/blog/article-1", MinDuration: 2, MaxDuration: 10, DropOffRate: 0.20},
			{EventName: "form_submit", URL: "/blog/article-1", MinDuration: 10, MaxDuration: 60, DropOffRate: 0.30},
		},
	},
	// Mobile App Install Journey
	{
		Name:    "Mobile App Install",
		Weight:  20,
		Country: "India",
		Browser: "Chrome",
		Device:  "Mobile",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/mobile", MinDuration: 5, MaxDuration: 20, DropOffRate: 0.30},
			{EventName: "button_click", URL: "/mobile", MinDuration: 2, MaxDuration: 10, DropOffRate: 0.25},
			{EventName: "app_install", URL: "/mobile", MinDuration: 30, MaxDuration: 180, DropOffRate: 0.60},
			{EventName: "app_open", URL: "/mobile", MinDuration: 5, MaxDuration: 3600, DropOffRate: 0.35},
		},
	},
	// Video Engagement
	{
		Name:    "Video Engagement",
		Weight:  15,
		Country: "Canada",
		Browser: "Safari",
		Device:  "Tablet",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/features", MinDuration: 5, MaxDuration: 20, DropOffRate: 0.20},
			{EventName: "video_play", URL: "/features", MinDuration: 10, MaxDuration: 60, DropOffRate: 0.45},
			{EventName: "button_click", URL: "/features", MinDuration: 2, MaxDuration: 10, DropOffRate: 0.35},
			{EventName: "signup", URL: "/signup", MinDuration: 30, MaxDuration: 120, DropOffRate: 0.50},
		},
	},
	// Support Journey
	{
		Name:    "Support Journey",
		Weight:  10,
		Country: "Australia",
		Browser: "Edge",
		Device:  "Desktop",
		Steps: []FunnelStep{
			{EventName: "page_view", URL: "/help", MinDuration: 10, MaxDuration: 60, DropOffRate: 0.25},
			{EventName: "search", URL: "/help", MinDuration: 5, MaxDuration: 30, DropOffRate: 0.40},
			{EventName: "page_view", URL: "/support", MinDuration: 20, MaxDuration: 180, DropOffRate: 0.50},
			{EventName: "button_click", URL: "/support", MinDuration: 2, MaxDuration: 10, DropOffRate: 0.30},
		},
	},
}

// User agent mapping
var userAgentMap = map[string]map[string]string{
	"Desktop": {
		"Chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Safari":  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
	},
	"Mobile": {
		"Chrome": "Mozilla/5.0 (Linux; Android 14; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		"Safari": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	},
	"Tablet": {
		"Chrome": "Mozilla/5.0 (Linux; Android 14; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Safari": "Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	},
}

// OS mapping
var osMap = map[string]string{
	"Desktop": "Windows",
	"Mobile":  "Android",
	"Tablet":  "iOS",
}

// IP ranges by country
var ipRangesByCountry = map[string]string{
	"United States": "203.0.113",
	"United Kingdom": "185.199.108",
	"Germany":        "151.101.1",
	"India":          "104.16.132",
	"Canada":         "198.51.100",
	"Australia":      "203.113.0",
	"Palestine":      "185.178.220",
}

// Referrers
var referrers = []string{
	"https://google.com/search?q=best+product",
	"https://facebook.com",
	"https://twitter.com",
	"https://linkedin.com",
	"https://reddit.com/r/technology",
	"direct",
	"email",
	"",
}

// selectFunnel selects a funnel template based on weights
func selectFunnel() FunnelTemplate {
	totalWeight := 0
	for _, ft := range funnelTemplates {
		totalWeight += ft.Weight
	}

	r := rand.Intn(totalWeight)
	currentWeight := 0

	for _, ft := range funnelTemplates {
		currentWeight += ft.Weight
		if r < currentWeight {
			return ft
		}
	}

	return funnelTemplates[0]
}

// generateUserJourney creates a sequence of events for one user following a funnel
func generateUserJourney(userID string, projectID string, baseTime time.Time) []Event {
	funnel := selectFunnel()

	sessionID := fmt.Sprintf("sess_%s_%d", userID, rand.Intn(100))
	currentTime := baseTime

	// Get user agent
	userAgent := userAgentMap[funnel.Device][funnel.Browser]
	if userAgent == "" {
		userAgent = userAgentMap["Desktop"]["Chrome"]
	}

	// Get OS
	os := osMap[funnel.Device]
	if os == "" {
		os = "Windows"
	}

	// Get IP range
	ipBase := ipRangesByCountry[funnel.Country]
	if ipBase == "" {
		ipBase = "192.168.1"
	}
	ip := fmt.Sprintf("%s.%d", ipBase, rand.Intn(255)+1)

	// Referrer
	referrer := referrers[rand.Intn(len(referrers))]

	var events []Event

	for i, step := range funnel.Steps {
		// Check drop-off
		if i > 0 && rand.Float64() < step.DropOffRate {
			// User dropped off at this step
			break
		}

		// Calculate time spent on this step
		duration := step.MinDuration + rand.Intn(step.MaxDuration-step.MinDuration+1)
		currentTime = currentTime.Add(time.Duration(duration) * time.Second)

		event := Event{
			Timestamp:       currentTime,
			EventName:       step.EventName,
			UserID:          userID,
			SessionID:       sessionID,
			SessionDuration: int(currentTime.Sub(baseTime).Seconds()),
			URL:             step.URL,
			Referrer:        referrer,
			UserAgent:       userAgent,
			IP:              ip,
			Country:         funnel.Country,
			Browser:         funnel.Browser,
			OS:              os,
			Device:          funnel.Device,
			IsBot:           false,
			ProjectID:       projectID,
		}

		events = append(events, event)

		// After first step, referrer becomes the previous URL
		referrer = step.URL
	}

	return events
}

// DBInserter handles database insertion
type DBInserter struct {
	db *sql.DB
}

func NewDBInserter(dbPath string) (*DBInserter, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DBInserter{db: db}, nil
}

func (di *DBInserter) Close() error {
	return di.db.Close()
}

func (di *DBInserter) InsertEvents(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := di.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration,
			url, referrer, user_agent, ip, country, browser, os, device, is_bot, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			event.SessionDuration,
			event.URL,
			event.Referrer,
			event.UserAgent,
			event.IP,
			event.Country,
			event.Browser,
			event.OS,
			event.Device,
			event.IsBot,
			event.ProjectID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// HTTPSender handles HTTP API requests
type HTTPSender struct {
	endpoint string
	client   *http.Client
}

func NewHTTPSender(endpoint string) *HTTPSender {
	return &HTTPSender{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (hs *HTTPSender) SendEvent(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", hs.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", event.UserAgent)

	resp, err := hs.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

func (hs *HTTPSender) SendEvents(events []Event) error {
	for _, event := range events {
		if err := hs.SendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	mode := flag.String("mode", "db", "Mode: 'db' or 'http'")
	numUsers := flag.Int("users", 10000, "Number of user journeys to generate")
	projectID := flag.String("project", "funnel_test", "Project ID")
	dbPath := flag.String("db", "../data/analytics.db", "Database path (for db mode)")
	endpoint := flag.String("endpoint", "http://localhost:8080/api/track", "API endpoint (for http mode)")
	daysBack := flag.Int("days", 30, "Generate data for the last N days")

	flag.Parse()

	log.Printf("🎯 Funnel Data Generator")
	log.Printf("📊 Mode: %s", *mode)
	log.Printf("👥 User Journeys: %d", *numUsers)
	log.Printf("🗓️  Time Range: Last %d days", *daysBack)
	log.Printf("📦 Project ID: %s", *projectID)

	// Time range: distribute users over the last N days
	baseTime := time.Now()
	timeRange := time.Duration(*daysBack) * 24 * time.Hour

	start := time.Now()
	totalEvents := 0

	switch *mode {
	case "db":
		log.Printf("💾 Database: %s", *dbPath)

		inserter, err := NewDBInserter(*dbPath)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer inserter.Close()

		batchSize := 100
		for i := 0; i < *numUsers; i++ {
			userID := fmt.Sprintf("funnel_user_%d", i+1)

			// Random time within the range
			userTime := baseTime.Add(-time.Duration(rand.Int63n(int64(timeRange))))

			events := generateUserJourney(userID, *projectID, userTime)
			totalEvents += len(events)

			if err := inserter.InsertEvents(events); err != nil {
				log.Printf("❌ Error inserting journey for %s: %v", userID, err)
				continue
			}

			if (i+1)%batchSize == 0 || i == *numUsers-1 {
				elapsed := time.Since(start)
				rate := float64(i+1) / elapsed.Seconds()
				log.Printf("📊 Progress: %d/%d journeys | %d events | %.0f journeys/sec",
					i+1, *numUsers, totalEvents, rate)
			}
		}

	case "http":
		log.Printf("🌐 Endpoint: %s", *endpoint)

		sender := NewHTTPSender(*endpoint)

		for i := 0; i < *numUsers; i++ {
			userID := fmt.Sprintf("funnel_user_%d", i+1)

			// Random time within the range
			userTime := baseTime.Add(-time.Duration(rand.Int63n(int64(timeRange))))

			events := generateUserJourney(userID, *projectID, userTime)
			totalEvents += len(events)

			if err := sender.SendEvents(events); err != nil {
				log.Printf("❌ Error sending journey for %s: %v", userID, err)
				continue
			}

			if (i+1)%100 == 0 || i == *numUsers-1 {
				elapsed := time.Since(start)
				rate := float64(i+1) / elapsed.Seconds()
				log.Printf("📊 Progress: %d/%d journeys | %d events | %.0f journeys/sec",
					i+1, *numUsers, totalEvents, rate)
			}
		}

	default:
		log.Fatal("Invalid mode. Use 'db' or 'http'")
	}

	duration := time.Since(start)
	log.Printf("\n✅ Funnel data generation completed!")
	log.Printf("📈 Total user journeys: %d", *numUsers)
	log.Printf("📈 Total events: %d", totalEvents)
	log.Printf("📊 Average events per journey: %.1f", float64(totalEvents)/float64(*numUsers))
	log.Printf("⏱️  Duration: %v", duration)

	// Print funnel statistics
	log.Printf("\n📊 Funnel Templates Used:")
	for _, ft := range funnelTemplates {
		log.Printf("   %s (%d steps, weight: %d)", ft.Name, len(ft.Steps), ft.Weight)
	}
}
