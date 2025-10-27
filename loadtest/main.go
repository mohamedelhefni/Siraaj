package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

// Event represents an analytics event for load testing
type Event struct {
	ID              uint64    `json:"id,omitempty"`
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

// Sample data for realistic events
var (
	eventNames = []string{
		"page_view", "button_click", "form_submit", "signup", "login", "logout",
		"purchase", "add_to_cart", "checkout_started", "payment_completed",
		"video_play", "video_pause", "search", "download", "share", "like",
		"comment", "follow", "unfollow", "profile_view", "settings_change",
		"notification_click", "email_open", "email_click", "app_install",
		"app_open", "feature_used", "error_occurred", "session_start", "session_end",
	}

	urls = []string{
		"/", "/home", "/about", "/contact", "/pricing", "/features", "/blog",
		"/login", "/signup", "/dashboard", "/profile", "/settings", "/help",
		"/product/123", "/product/456", "/product/789", "/category/electronics",
		"/category/clothing", "/category/books", "/search?q=laptop", "/cart",
		"/checkout", "/payment", "/confirmation", "/account", "/orders", "/support",
	}

	referrers = []string{
		"", "https://google.com", "https://facebook.com", "https://twitter.com",
		"https://linkedin.com", "https://reddit.com", "https://youtube.com",
		"https://github.com", "https://stackoverflow.com", "https://medium.com",
		"https://dev.to", "https://hackernews.com", "direct", "email",
	}

	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	}

	botUserAgents = []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
		"Twitterbot/1.0",
	}

	countries = []string{
		"United States", "Canada", "United Kingdom", "Germany", "France", "Spain",
		"Italy", "Netherlands", "Sweden", "Norway", "Denmark", "Finland",
		"Australia", "New Zealand", "Japan", "South Korea", "Singapore", "India",
		"Brazil", "Mexico", "Argentina", "Chile", "Russia", "China", "Palestine",
	}

	browsers = []string{
		"Chrome", "Safari", "Firefox", "Edge", "Opera", "Brave",
	}

	operatingSystems = []string{
		"Windows", "MacOS", "Linux", "iOS", "Android", "ChromeOS",
	}

	devices = []string{
		"Desktop", "Mobile", "Tablet",
	}

	ipRanges = []string{
		"192.168.1", "10.0.0", "172.16.0", "203.0.113", "198.51.100",
		"203.113.0", "185.199.108", "140.82.112", "151.101.1", "104.16.132",
	}
)

// GenerateRandomEvent creates a realistic random event
func GenerateRandomEvent(baseTime time.Time, userPool []string, projectID string) Event {
	// Random timestamp within the last 30 days
	hoursBack := rand.Intn(30 * 24)
	timestamp := baseTime.Add(-time.Duration(hoursBack) * time.Hour)
	timestamp = timestamp.Add(time.Duration(rand.Intn(3600)) * time.Second)

	// Select random user from pool
	userID := userPool[rand.Intn(len(userPool))]

	// Generate session ID
	sessionID := fmt.Sprintf("sess_%s_%d", userID, rand.Intn(10))

	// Generate session duration (0-3600 seconds, 1 hour max)
	sessionDuration := rand.Intn(3600)

	eventName := eventNames[rand.Intn(len(eventNames))]
	url := urls[rand.Intn(len(urls))]
	referrer := referrers[rand.Intn(len(referrers))]

	// 20% chance of bot
	isBot := rand.Float32() < 0.2
	var userAgent string
	if isBot {
		userAgent = botUserAgents[rand.Intn(len(botUserAgents))]
	} else {
		userAgent = userAgents[rand.Intn(len(userAgents))]
	}

	country := countries[rand.Intn(len(countries))]
	browser := browsers[rand.Intn(len(browsers))]
	os := operatingSystems[rand.Intn(len(operatingSystems))]
	device := devices[rand.Intn(len(devices))]

	// Generate realistic IP
	ipBase := ipRanges[rand.Intn(len(ipRanges))]
	ip := fmt.Sprintf("%s.%d", ipBase, rand.Intn(255)+1)

	return Event{
		Timestamp:       timestamp,
		EventName:       eventName,
		UserID:          userID,
		SessionID:       sessionID,
		SessionDuration: sessionDuration,
		URL:             url,
		Referrer:        referrer,
		UserAgent:       userAgent,
		IP:              ip,
		Country:         country,
		Browser:         browser,
		OS:              os,
		Device:          device,
		IsBot:           isBot,
		ProjectID:       projectID,
	}
}

// ========== DATABASE MODE ==========

type DBLoadTester struct {
	db *sql.DB
}

func NewDBLoadTester(dbPath string) (*DBLoadTester, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DBLoadTester{db: db}, nil
}

func (lt *DBLoadTester) Close() error {
	return lt.db.Close()
}

func (lt *DBLoadTester) InsertEventsBatch(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := lt.db.Begin()
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

func (lt *DBLoadTester) RunLoadTest(totalEvents int, batchSize int, numUsers int, projectID string) error {
	log.Printf("🚀 Starting DB load test: %d events, batch size: %d, users: %d", totalEvents, batchSize, numUsers)

	userPool := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userPool[i] = fmt.Sprintf("user_%d", i+1)
	}

	baseTime := time.Now()
	totalBatches := (totalEvents + batchSize - 1) / batchSize

	start := time.Now()

	for batch := 0; batch < totalBatches; batch++ {
		batchStart := time.Now()

		eventsInBatch := batchSize
		if batch == totalBatches-1 {
			eventsInBatch = totalEvents - (batch * batchSize)
		}

		events := make([]Event, eventsInBatch)
		for i := 0; i < eventsInBatch; i++ {
			events[i] = GenerateRandomEvent(baseTime, userPool, projectID)
		}

		if err := lt.InsertEventsBatch(events); err != nil {
			return fmt.Errorf("error inserting batch %d: %v", batch+1, err)
		}

		batchDuration := time.Since(batchStart)
		totalInserted := (batch + 1) * batchSize
		if totalInserted > totalEvents {
			totalInserted = totalEvents
		}

		if batch%10 == 0 || batch == totalBatches-1 {
			elapsed := time.Since(start)
			rate := float64(totalInserted) / elapsed.Seconds()

			log.Printf("📊 Batch %d/%d | Events: %d/%d | Rate: %.0f events/sec | Batch time: %v",
				batch+1, totalBatches, totalInserted, totalEvents, rate, batchDuration)
		}
	}

	duration := time.Since(start)
	rate := float64(totalEvents) / duration.Seconds()

	log.Printf("✅ DB load test completed!")
	log.Printf("📈 Total events: %d", totalEvents)
	log.Printf("⏱️  Total time: %v", duration)
	log.Printf("🚄 Average rate: %.0f events/sec", rate)

	return nil
}

// ========== HTTP MODE ==========

type HTTPLoadTester struct {
	endpoint string
	client   *http.Client
}

func NewHTTPLoadTester(endpoint string) *HTTPLoadTester {
	// High-performance HTTP client
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		WriteBufferSize:     32 * 1024,
		ReadBufferSize:      32 * 1024,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return &HTTPLoadTester{
		endpoint: endpoint,
		client:   client,
	}
}

func (ht *HTTPLoadTester) SendEvent(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ht.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ht.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Discard response body
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (ht *HTTPLoadTester) RunLoadTest(totalEvents int, workers int, numUsers int, projectID string) error {
	log.Printf("🚀 Starting HTTP load test: %d events, %d workers, %d users", totalEvents, workers, numUsers)

	userPool := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userPool[i] = fmt.Sprintf("user_%d", i+1)
	}

	baseTime := time.Now()
	start := time.Now()

	var successCount atomic.Int64
	var errorCount atomic.Int64
	var wg sync.WaitGroup

	// Channel for events
	eventChan := make(chan Event, workers*2)

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range eventChan {
				if err := ht.SendEvent(event); err != nil {
					errorCount.Add(1)
				} else {
					successCount.Add(1)
				}
			}
		}()
	}

	// Progress reporter
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(start)
				success := successCount.Load()
				errors := errorCount.Load()
				rate := float64(success) / elapsed.Seconds()
				log.Printf("📊 Progress: %d sent, %d errors | Rate: %.0f events/sec", success, errors, rate)
			case <-done:
				return
			}
		}
	}()

	// Generate and send events
	for i := 0; i < totalEvents; i++ {
		event := GenerateRandomEvent(baseTime, userPool, projectID)
		eventChan <- event
	}

	close(eventChan)
	wg.Wait()
	close(done)

	duration := time.Since(start)
	success := successCount.Load()
	errors := errorCount.Load()
	rate := float64(success) / duration.Seconds()

	log.Printf("✅ HTTP load test completed!")
	log.Printf("📈 Total events sent: %d", success)
	log.Printf("❌ Errors: %d", errors)
	log.Printf("⏱️  Total time: %v", duration)
	log.Printf("🚄 Average rate: %.0f events/sec", rate)

	return nil
}

// ========== CSV GENERATOR ==========

type CSVGenerator struct {
	filepath string
}

func NewCSVGenerator(filepath string) *CSVGenerator {
	return &CSVGenerator{filepath: filepath}
}

func (cg *CSVGenerator) GenerateCSV(totalEvents int, numUsers int, projectID string) error {
	log.Printf("📝 Generating CSV file: %s with %d events", cg.filepath, totalEvents)
	
	file, err := os.Create(cg.filepath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{
		"timestamp", "event_name", "user_id", "session_id", "session_duration",
		"url", "referrer", "user_agent", "ip", "country", "browser", "os", "device", "is_bot", "project_id",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	userPool := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userPool[i] = fmt.Sprintf("user_%d", i+1)
	}

	baseTime := time.Now()
	start := time.Now()
	
	for i := 0; i < totalEvents; i++ {
		event := GenerateRandomEvent(baseTime, userPool, projectID)
		
		record := []string{
			event.Timestamp.Format(time.RFC3339),
			event.EventName,
			event.UserID,
			event.SessionID,
			fmt.Sprintf("%d", event.SessionDuration),
			event.URL,
			event.Referrer,
			event.UserAgent,
			event.IP,
			event.Country,
			event.Browser,
			event.OS,
			event.Device,
			fmt.Sprintf("%t", event.IsBot),
			event.ProjectID,
		}
		
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
		
		if (i+1)%100000 == 0 {
			elapsed := time.Since(start)
			rate := float64(i+1) / elapsed.Seconds()
			log.Printf("📊 Progress: %d/%d events written (%.0f events/sec)", i+1, totalEvents, rate)
		}
	}

	duration := time.Since(start)
	rate := float64(totalEvents) / duration.Seconds()
	
	log.Printf("✅ CSV generation completed!")
	log.Printf("📈 Total events: %d", totalEvents)
	log.Printf("⏱️  Total time: %v", duration)
	log.Printf("🚄 Average rate: %.0f events/sec", rate)
	
	return nil
}

func (cg *CSVGenerator) ImportToDatabase(dbPath string) error {
	log.Printf("📥 Importing CSV file %s to database %s", cg.filepath, dbPath)
	
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	start := time.Now()
	
	// Import CSV using DuckDB's INSERT INTO ... SELECT FROM read_csv_auto
	// Use nextval('id_sequence') for auto-incrementing IDs
	query := fmt.Sprintf(`
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration, url, referrer, user_agent, ip, country, browser, os, device, is_bot, project_id)
		SELECT nextval('id_sequence'), timestamp::TIMESTAMP, event_name, user_id, session_id, session_duration::INTEGER, url, referrer, user_agent, ip, country, browser, os, device, is_bot::BOOLEAN, project_id
		FROM read_csv_auto('%s', header=true)
	`, cg.filepath)
	
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to import CSV: %w", err)
	}
	
	duration := time.Since(start)
	
	log.Printf("✅ CSV import completed in %v", duration)
	
	return nil
}

// ========== MAIN ==========

func main() {
	mode := flag.String("mode", "db", "Load test mode: 'db', 'http', or 'csv'")
	events := flag.Int("events", 1000000, "Total number of events to generate")
	batchSize := flag.Int("batch", 1000, "Batch size for DB mode")
	workers := flag.Int("workers", 50, "Number of concurrent workers for HTTP mode")
	users := flag.Int("users", 10000, "Number of unique users to simulate")
	projectID := flag.String("project", "test_project", "Project ID for events")
	dbPath := flag.String("db", "../data/analytics.db", "Database path for DB mode")
	endpoint := flag.String("endpoint", "http://localhost:8080/api/events", "API endpoint for HTTP mode")
	csvPath := flag.String("csv", "../data/loadtest.csv", "CSV file path for CSV mode")

	flag.Parse()

	log.Printf("🔧 Configuration:")
	log.Printf("  Mode: %s", *mode)
	log.Printf("  Events: %d", *events)
	log.Printf("  Users: %d", *users)
	log.Printf("  Project ID: %s", *projectID)

	switch *mode {
	case "db":
		log.Printf("  Database: %s", *dbPath)
		log.Printf("  Batch Size: %d", *batchSize)

		lt, err := NewDBLoadTester(*dbPath)
		if err != nil {
			log.Fatal("Failed to create DB load tester:", err)
		}
		defer lt.Close()

		if err := lt.RunLoadTest(*events, *batchSize, *users, *projectID); err != nil {
			log.Fatal("DB load test failed:", err)
		}

	case "http":
		log.Printf("  Endpoint: %s", *endpoint)
		log.Printf("  Workers: %d", *workers)

		ht := NewHTTPLoadTester(*endpoint)
		if err := ht.RunLoadTest(*events, *workers, *users, *projectID); err != nil {
			log.Fatal("HTTP load test failed:", err)
		}

	case "csv":
		log.Printf("  CSV Path: %s", *csvPath)
		log.Printf("  Database: %s", *dbPath)
		
		cg := NewCSVGenerator(*csvPath)
		
		// Generate CSV
		if err := cg.GenerateCSV(*events, *users, *projectID); err != nil {
			log.Fatal("CSV generation failed:", err)
		}
		
		// Import to database
		if err := cg.ImportToDatabase(*dbPath); err != nil {
			log.Fatal("CSV import failed:", err)
		}

	default:
		log.Fatal("Invalid mode. Use 'db', 'http', or 'csv'")
	}
}
