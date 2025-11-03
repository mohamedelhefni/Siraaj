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
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
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
	Channel         string    `json:"channel"`
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
		// Paid channels
		"https://google.com/ads", "https://facebook.com/ads", "https://twitter.com/ads",
		"https://linkedin.com/ads", "https://instagram.com/ads",
		// Organic search
		"https://www.google.com/search", "https://www.bing.com/search", "https://search.yahoo.com",
		// Social
		"https://t.co", "https://www.facebook.com", "https://www.linkedin.com",
		"https://www.instagram.com", "https://www.tiktok.com",
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

// DetectChannel determines the traffic channel based on referrer and URL
func DetectChannel(referrer, url string) string {
	// Priority 1: Paid channels (utm_medium or utm_source contains paid/cpc/ppc)
	if containsAny(url, []string{"utm_medium=cpc", "utm_medium=ppc", "utm_medium=paid", "utm_source=paid"}) ||
		containsAny(referrer, []string{"/ads", "adwords", "googleads", "facebook.com/ads"}) {
		return "Paid"
	}

	// Priority 2: Direct traffic (no referrer or same domain)
	if referrer == "" || referrer == "direct" {
		return "Direct"
	}

	// Priority 3: Social media
	socialDomains := []string{
		"facebook.com", "twitter.com", "linkedin.com", "instagram.com",
		"tiktok.com", "pinterest.com", "reddit.com", "youtube.com",
		"snapchat.com", "whatsapp.com", "telegram.org", "t.co",
	}
	if containsAnyDomain(referrer, socialDomains) {
		return "Social"
	}

	// Priority 4: Organic search
	searchEngines := []string{
		"google.com/search", "bing.com/search", "yahoo.com/search",
		"duckduckgo.com", "baidu.com", "yandex.com", "ask.com",
	}
	if containsAny(referrer, searchEngines) {
		return "Organic"
	}

	// Priority 5: Referral (all other external sources)
	return "Referral"
}

// containsAny checks if the text contains any of the substrings
func containsAny(text string, substrings []string) bool {
	for _, substr := range substrings {
		if len(text) >= len(substr) && indexOfSubstring(text, substr) >= 0 {
			return true
		}
	}
	return false
}

// containsAnyDomain checks if the referrer URL contains any of the domains
func containsAnyDomain(referrer string, domains []string) bool {
	for _, domain := range domains {
		if indexOfSubstring(referrer, domain) >= 0 {
			return true
		}
	}
	return false
}

// indexOfSubstring finds the index of a substring (case-insensitive)
func indexOfSubstring(s, substr string) int {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return i
		}
	}
	return -1
}

// toLower converts a string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

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

	// Detect channel based on referrer and URL
	channel := DetectChannel(referrer, url)

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
		Channel:         channel,
	}
}

// ========== DATABASE MODE ==========

type DBLoadTester struct {
	db *sql.DB
}

func NewDBLoadTester(dsn string) (*DBLoadTester, error) {
	// Add ClickHouse performance settings to DSN
	// Enable async_insert for extreme speed
	if len(dsn) > 0 && dsn[len(dsn)-1] != '&' && dsn[len(dsn)-1] != '?' {
		if containsAny(dsn, []string{"?"}) {
			dsn += "&"
		} else {
			dsn += "?"
		}
	}
	dsn += "async_insert=1&wait_for_async_insert=0&async_insert_max_data_size=10000000&async_insert_busy_timeout_ms=1000"

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool for maximum throughput
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Hour)

	return &DBLoadTester{db: db}, nil
}

func (lt *DBLoadTester) Close() error {
	return lt.db.Close()
}

func (lt *DBLoadTester) InsertEventsBatch(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	// ClickHouse ultra-fast batch insert using async_insert
	// Pre-allocate with exact capacity for maximum performance
	const columnsPerRow = 17
	query := `INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration,
		url, referrer, user_agent, ip, country, browser, os, device, is_bot, project_id, channel) VALUES `

	values := make([]interface{}, 0, len(events)*columnsPerRow)

	// Pre-build placeholder string for reuse (much faster)
	placeholder := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// Build query string efficiently
	queryBuilder := make([]byte, 0, len(query)+len(events)*(len(placeholder)+2))
	queryBuilder = append(queryBuilder, query...)
	queryBuilder = append(queryBuilder, placeholder...)

	baseNano := uint64(time.Now().UnixNano())

	for i, event := range events {
		if i > 0 {
			queryBuilder = append(queryBuilder, ',', ' ')
			queryBuilder = append(queryBuilder, placeholder...)
		}

		eventID := baseNano + uint64(i)

		values = append(values,
			eventID,
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
			event.Channel,
		)
	}

	// Execute with async_insert - ClickHouse buffers and writes asynchronously
	_, err := lt.db.Exec(string(queryBuilder), values...)
	return err
}

func (lt *DBLoadTester) RunLoadTest(totalEvents int, batchSize int, numUsers int, projectID string) error {
	log.Printf("🚀 Starting EXTREME SPEED ClickHouse load test!")
	log.Printf("   Events: %d | Batch size: %d | Users: %d", totalEvents, batchSize, numUsers)
	log.Printf("   Using async_insert for maximum throughput 🔥")

	// Pre-generate user pool
	userPool := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userPool[i] = fmt.Sprintf("user_%d", i+1)
	}

	baseTime := time.Now()
	totalBatches := (totalEvents + batchSize - 1) / batchSize

	// Use worker pool for parallel batch generation and insertion
	numWorkers := 10 // Parallel workers for extreme speed
	batchChan := make(chan int, numWorkers*2)
	errorsChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	// Progress tracking
	var completed atomic.Int64
	start := time.Now()

	// Progress reporter
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		lastCount := int64(0)
		for {
			select {
			case <-ticker.C:
				current := completed.Load()
				elapsed := time.Since(start)
				rate := float64(current) / elapsed.Seconds()
				instantRate := float64(current - lastCount) // events in last second
				lastCount = current

				log.Printf("⚡ Progress: %d/%d | Speed: %.0f events/sec | Instant: %.0f/sec",
					current, totalEvents, rate, instantRate)
			case <-done:
				return
			}
		}
	}()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batchNum := range batchChan {
				eventsInBatch := batchSize
				if batchNum == totalBatches-1 {
					eventsInBatch = totalEvents - (batchNum * batchSize)
				}

				// Generate events for this batch
				events := make([]Event, eventsInBatch)
				for i := 0; i < eventsInBatch; i++ {
					events[i] = GenerateRandomEvent(baseTime, userPool, projectID)
				}

				// Insert batch
				if err := lt.InsertEventsBatch(events); err != nil {
					errorsChan <- fmt.Errorf("worker %d, batch %d: %v", workerID, batchNum+1, err)
					return
				}

				completed.Add(int64(eventsInBatch))
			}
		}(i)
	}

	// Distribute batches to workers
	go func() {
		for batch := 0; batch < totalBatches; batch++ {
			batchChan <- batch
		}
		close(batchChan)
	}()

	// Wait for completion
	wg.Wait()
	close(done)
	close(errorsChan)

	// Check for errors
	if len(errorsChan) > 0 {
		return <-errorsChan
	}

	duration := time.Since(start)
	rate := float64(totalEvents) / duration.Seconds()

	log.Printf("\n✅ EXTREME SPEED load test completed!")
	log.Printf("📈 Total events inserted: %d", totalEvents)
	log.Printf("⏱️  Total time: %v", duration)
	log.Printf("🚄 Average rate: %.0f events/sec", rate)
	log.Printf("💾 Data size: ~%.2f MB", float64(totalEvents*500)/1024/1024) // ~500 bytes per event

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
		"url", "referrer", "user_agent", "ip", "country", "browser", "os", "device", "is_bot", "project_id", "channel",
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
			event.Channel,
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

func (cg *CSVGenerator) ImportToDatabase(dsn string) error {
	log.Printf("📥 Importing CSV file %s to ClickHouse", cg.filepath)

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	start := time.Now()

	// Read CSV and insert to ClickHouse
	file, err := os.Open(cg.filepath)
	if err != nil {
		return fmt.Errorf("failed to open CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Read and batch insert records (ClickHouse doesn't use traditional transactions)
	batchSize := 10000
	batch := make([][]string, 0, batchSize)
	counter := uint64(0)

	for {
		record, err := reader.Read()
		if err != nil {
			// Process final batch
			if len(batch) > 0 {
				if err := insertCSVBatch(db, batch, &counter); err != nil {
					return err
				}
			}
			break
		}

		batch = append(batch, record)

		if len(batch) >= batchSize {
			if err := insertCSVBatch(db, batch, &counter); err != nil {
				return err
			}
			batch = batch[:0] // Reset batch
			log.Printf("📊 Imported %d rows...", counter)
		}
	}

	duration := time.Since(start)
	log.Printf("✅ CSV import completed in %v", duration)

	return nil
}

// insertCSVBatch inserts a batch of CSV records into ClickHouse
func insertCSVBatch(db *sql.DB, records [][]string, counter *uint64) error {
	if len(records) == 0 {
		return nil
	}

	query := `INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration,
		url, referrer, user_agent, ip, country, browser, os, device, is_bot, project_id, channel) VALUES `

	values := make([]interface{}, 0, len(records)*17)
	valuePlaceholders := make([]string, 0, len(records))

	for _, record := range records {
		*counter++
		timestamp, _ := time.Parse(time.RFC3339, record[0])
		sessionDuration, _ := strconv.Atoi(record[4])
		isBot, _ := strconv.ParseBool(record[13])
		eventID := uint64(time.Now().UnixNano()) + *counter

		valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		values = append(values,
			eventID,
			timestamp,
			record[1], // event_name
			record[2], // user_id
			record[3], // session_id
			sessionDuration,
			record[5],  // url
			record[6],  // referrer
			record[7],  // user_agent
			record[8],  // ip
			record[9],  // country
			record[10], // browser
			record[11], // os
			record[12], // device
			isBot,
			record[14], // project_id
			record[15], // channel
		)
	}

	query += valuePlaceholders[0]
	for i := 1; i < len(valuePlaceholders); i++ {
		query += ", " + valuePlaceholders[i]
	}

	_, err := db.Exec(query, values...)
	return err
}

// ImportToParquet is removed as ClickHouse handles storage natively
// CSV data should be imported directly using ImportToDatabase

// csvRecordToEvent converts CSV record to Event struct (for HTTP mode)
func csvRecordToEvent(record []string, currentID *uint64) (Event, error) {
	*currentID++

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, record[0])
	if err != nil {
		timestamp = time.Now()
	}

	// Parse session duration
	sessionDuration, _ := strconv.Atoi(record[4])

	// Parse is_bot
	isBot, _ := strconv.ParseBool(record[13])

	return Event{
		ID:              *currentID,
		Timestamp:       timestamp,
		EventName:       record[1],
		UserID:          record[2],
		SessionID:       record[3],
		SessionDuration: sessionDuration,
		URL:             record[5],
		Referrer:        record[6],
		UserAgent:       record[7],
		IP:              record[8],
		Country:         record[9],
		Browser:         record[10],
		OS:              record[11],
		Device:          record[12],
		IsBot:           isBot,
		ProjectID:       record[14],
		Channel:         record[15],
	}, nil
}

// ========== MAIN ==========

func main() {
	mode := flag.String("mode", "db", "Load test mode: 'db', 'http', or 'csv'")
	events := flag.Int("events", 10000000, "Total number of events to generate (default: 10M for speed test)")
	batchSize := flag.Int("batch", 50000, "Batch size for DB mode (default: 50K for ClickHouse async_insert)")
	workers := flag.Int("workers", 50, "Number of concurrent workers for HTTP mode")
	users := flag.Int("users", 100000, "Number of unique users to simulate (default: 100K)")
	projectID := flag.String("project", "test_project", "Project ID for events")
	dsn := flag.String("dsn", "clickhouse://localhost:9000/siraaj?username=default&password=", "ClickHouse DSN for DB mode")
	endpoint := flag.String("endpoint", "http://localhost:8080/api/events", "API endpoint for HTTP mode")
	csvPath := flag.String("csv", "../data/loadtest.csv", "CSV file path for CSV mode")

	flag.Parse()

	log.Printf("🔧 Configuration:")
	log.Printf("  Mode: %s", *mode)
	log.Printf("  Events: %d", *events)
	log.Printf("  Users: %d", *users)
	log.Printf("  Project ID: %s", *projectID)

	if *mode == "db" {
		log.Printf("\n💡 TIP: For maximum ClickHouse performance:")
		log.Printf("   - Larger batches = faster (try -batch=100000)")
		log.Printf("   - async_insert is automatically enabled")
		log.Printf("   - 10 parallel workers are used")
		log.Printf("   - Expected: 200K-500K events/sec on modern hardware\n")
	}

	switch *mode {
	case "db":
		log.Printf("  ClickHouse DSN: %s", *dsn)
		log.Printf("  Batch Size: %d", *batchSize)

		lt, err := NewDBLoadTester(*dsn)
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

		cg := NewCSVGenerator(*csvPath)

		// Generate CSV
		if err := cg.GenerateCSV(*events, *users, *projectID); err != nil {
			log.Fatal("CSV generation failed:", err)
		}

		// Import to ClickHouse
		log.Println("📥 Import CSV to ClickHouse using: -mode=db with ImportToDatabase")
		log.Printf("Example: go run main.go -mode=db -dsn=\"%s\" -csv=%s", *dsn, *csvPath)

	default:
		log.Fatal("Invalid mode. Use 'db', 'http', or 'csv'")
	}
}
