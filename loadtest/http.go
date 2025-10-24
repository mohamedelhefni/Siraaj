package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Event represents an analytics event for HTTP tracking
type Event struct {
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
	Properties string    `json:"properties"` // JSON-encoded properties
}

// HTTPLoadTester handles HTTP-based load testing
type HTTPLoadTester struct {
	baseURL    string
	httpClient *http.Client
	stats      *LoadTestStats
}

// LoadTestStats tracks performance metrics
type LoadTestStats struct {
	totalRequests  int64
	successfulReqs int64
	failedReqs     int64
	totalLatency   int64 // microseconds
	minLatency     int64
	maxLatency     int64
	startTime      time.Time
	endTime        time.Time
	responseCodes  map[int]int64
	mutex          sync.RWMutex
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

	countries = []string{
		"United States", "Canada", "United Kingdom", "Germany", "France", "Spain",
		"Italy", "Netherlands", "Sweden", "Norway", "Denmark", "Finland",
		"Australia", "New Zealand", "Japan", "South Korea", "Singapore", "India",
		"Brazil", "Mexico", "Argentina", "Chile", "Russia", "China", "Israel",
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

func NewHTTPLoadTester(baseURL string) *HTTPLoadTester {
	return &HTTPLoadTester{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		stats: &LoadTestStats{
			responseCodes: make(map[int]int64),
			minLatency:    1000000, // Start with a high value
		},
	}
}

// GenerateRandomEvent creates a realistic random event
func (hlt *HTTPLoadTester) GenerateRandomEvent(userPool []string) Event {
	// Random timestamp within the last hour for real-time simulation
	minutesBack := rand.Intn(60) // last hour
	timestamp := time.Now().Add(-time.Duration(minutesBack) * time.Minute)

	// Select random user from pool
	userID := userPool[rand.Intn(len(userPool))]
	sessionID := fmt.Sprintf("sess_%s_%d", userID, rand.Intn(10))

	eventName := eventNames[rand.Intn(len(eventNames))]
	url := urls[rand.Intn(len(urls))]
	referrer := referrers[rand.Intn(len(referrers))]
	userAgent := userAgents[rand.Intn(len(userAgents))]
	country := countries[rand.Intn(len(countries))]
	browser := browsers[rand.Intn(len(browsers))]
	os := operatingSystems[rand.Intn(len(operatingSystems))]
	device := devices[rand.Intn(len(devices))]

	// Generate realistic IP
	ipBase := ipRanges[rand.Intn(len(ipRanges))]
	ip := fmt.Sprintf("%s.%d", ipBase, rand.Intn(255)+1)

	// Generate properties based on event type
	properties := hlt.generateProperties(eventName)
	encodedProps, _ := json.Marshal(properties)

	return Event{
		Timestamp:  timestamp,
		EventName:  eventName,
		UserID:     userID,
		SessionID:  sessionID,
		URL:        url,
		Referrer:   referrer,
		UserAgent:  userAgent,
		IP:         ip,
		Country:    country,
		Browser:    browser,
		OS:         os,
		Device:     device,
		Properties: string(encodedProps),
	}
}

// generateProperties creates realistic properties based on event type
func (hlt *HTTPLoadTester) generateProperties(eventName string) map[string]interface{} {
	props := make(map[string]interface{})

	switch eventName {
	case "purchase", "payment_completed":
		props["amount"] = rand.Float64()*1000 + 10
		props["currency"] = "USD"
		props["product_id"] = fmt.Sprintf("PROD-%d", rand.Intn(1000))
		props["category"] = []string{"electronics", "clothing", "books", "software"}[rand.Intn(4)]
		props["quantity"] = rand.Intn(5) + 1

	case "signup", "login":
		props["method"] = []string{"email", "google", "facebook", "github"}[rand.Intn(4)]
		props["source"] = []string{"organic", "paid", "referral", "direct"}[rand.Intn(4)]

	case "button_click":
		props["button_id"] = []string{"cta_hero", "nav_signup", "footer_contact", "product_buy"}[rand.Intn(4)]
		props["position"] = []string{"header", "hero", "sidebar", "footer"}[rand.Intn(4)]

	case "search":
		props["query"] = []string{"laptop", "phone", "headphones", "camera", "tablet"}[rand.Intn(5)]
		props["results_count"] = rand.Intn(100)

	case "video_play", "video_pause":
		props["video_id"] = fmt.Sprintf("VID-%d", rand.Intn(500))
		props["duration"] = rand.Intn(3600)
		props["position"] = rand.Intn(1800)

	case "error_occurred":
		props["error_type"] = []string{"network", "validation", "timeout", "404", "500"}[rand.Intn(5)]
		props["severity"] = []string{"low", "medium", "high", "critical"}[rand.Intn(4)]

	case "add_to_cart":
		props["product_id"] = fmt.Sprintf("PROD-%d", rand.Intn(1000))
		props["price"] = rand.Float64()*500 + 5
		props["category"] = []string{"electronics", "clothing", "books"}[rand.Intn(3)]
	}

	// Add common properties
	props["timestamp"] = time.Now().Unix()
	props["page_load_time"] = rand.Intn(5000) + 500

	return props
}

// SendEvent sends a single event to the tracking endpoint
func (hlt *HTTPLoadTester) SendEvent(event Event) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	start := time.Now()

	req, err := http.NewRequest("POST", hlt.baseURL+"/api/track", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", event.UserAgent)

	resp, err := hlt.httpClient.Do(req)
	latency := time.Since(start).Microseconds()

	// Update stats
	hlt.updateStats(latency, resp, err)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// updateStats safely updates the performance statistics
func (hlt *HTTPLoadTester) updateStats(latency int64, resp *http.Response, err error) {
	hlt.stats.mutex.Lock()
	defer hlt.stats.mutex.Unlock()

	atomic.AddInt64(&hlt.stats.totalRequests, 1)
	atomic.AddInt64(&hlt.stats.totalLatency, latency)

	if err != nil {
		atomic.AddInt64(&hlt.stats.failedReqs, 1)
		if resp != nil {
			hlt.stats.responseCodes[resp.StatusCode]++
		}
	} else {
		atomic.AddInt64(&hlt.stats.successfulReqs, 1)
		if resp != nil {
			hlt.stats.responseCodes[resp.StatusCode]++
		}
	}

	// Update min/max latency
	if latency < hlt.stats.minLatency {
		hlt.stats.minLatency = latency
	}
	if latency > hlt.stats.maxLatency {
		hlt.stats.maxLatency = latency
	}
}

// RunConcurrentLoadTest runs a concurrent HTTP load test
func (hlt *HTTPLoadTester) RunConcurrentLoadTest(totalRequests int, concurrency int, numUsers int, duration time.Duration) error {
	log.Printf("🚀 Starting HTTP load test")
	log.Printf("📊 Target: %s/api/track", hlt.baseURL)
	log.Printf("📈 Total requests: %d", totalRequests)
	log.Printf("🔄 Concurrency: %d", concurrency)
	log.Printf("👥 Users: %d", numUsers)
	log.Printf("⏱️  Duration: %v", duration)

	// Test server connectivity first
	resp, err := hlt.httpClient.Get(hlt.baseURL + "/api/health")
	if err != nil {
		return fmt.Errorf("❌ Cannot connect to server: %v", err)
	}
	resp.Body.Close()
	log.Printf("✅ Server connectivity verified")

	// Generate user pool
	userPool := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userPool[i] = fmt.Sprintf("loadtest_user_%d", i+1)
	}

	hlt.stats.startTime = time.Now()

	// Create channels for work distribution
	eventChan := make(chan Event, concurrency*2)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for event := range eventChan {
				if err := hlt.SendEvent(event); err != nil {
					log.Printf("Worker %d error: %v", workerID, err)
				}
			}
		}(i)
	}

	// Generate and send events
	go func() {
		defer close(eventChan)

		requestCount := 0
		timeout := time.After(duration)
		ticker := time.NewTicker(time.Second / time.Duration(totalRequests/int(duration.Seconds())))
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				log.Printf("⏰ Duration limit reached")
				return
			case <-ticker.C:
				if requestCount >= totalRequests {
					log.Printf("🎯 Request limit reached")
					return
				}

				event := hlt.GenerateRandomEvent(userPool)
				select {
				case eventChan <- event:
					requestCount++
				default:
					// Channel full, skip this event
				}
			}
		}
	}()

	// Progress reporting
	go func() {
		progressTicker := time.NewTicker(5 * time.Second)
		defer progressTicker.Stop()

		for {
			select {
			case <-progressTicker.C:
				hlt.printProgress()
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()
	hlt.stats.endTime = time.Now()

	log.Printf("✅ Load test completed!")
	hlt.printFinalStats()

	return nil
}

// printProgress displays current progress
func (hlt *HTTPLoadTester) printProgress() {
	hlt.stats.mutex.RLock()
	defer hlt.stats.mutex.RUnlock()

	total := atomic.LoadInt64(&hlt.stats.totalRequests)
	successful := atomic.LoadInt64(&hlt.stats.successfulReqs)
	failed := atomic.LoadInt64(&hlt.stats.failedReqs)
	avgLatency := int64(0)

	if total > 0 {
		avgLatency = atomic.LoadInt64(&hlt.stats.totalLatency) / total
	}

	elapsed := time.Since(hlt.stats.startTime)
	rps := float64(total) / elapsed.Seconds()

	log.Printf("📊 Progress: %d requests | ✅ %d success | ❌ %d failed | 🚄 %.1f req/s | ⏱️ %.1fms avg",
		total, successful, failed, rps, float64(avgLatency)/1000.0)
}

// printFinalStats displays comprehensive final statistics
func (hlt *HTTPLoadTester) printFinalStats() {
	hlt.stats.mutex.RLock()
	defer hlt.stats.mutex.RUnlock()

	total := atomic.LoadInt64(&hlt.stats.totalRequests)
	successful := atomic.LoadInt64(&hlt.stats.successfulReqs)
	failed := atomic.LoadInt64(&hlt.stats.failedReqs)
	totalLatency := atomic.LoadInt64(&hlt.stats.totalLatency)

	duration := hlt.stats.endTime.Sub(hlt.stats.startTime)
	rps := float64(total) / duration.Seconds()
	avgLatency := float64(0)

	if total > 0 {
		avgLatency = float64(totalLatency) / float64(total) / 1000.0 // Convert to milliseconds
	}

	successRate := float64(successful) / float64(total) * 100

	fmt.Println("\n" + "=")
	fmt.Println("📊 FINAL LOAD TEST RESULTS")
	fmt.Println("=")
	fmt.Printf("📈 Total Requests:     %d\n", total)
	fmt.Printf("✅ Successful:         %d (%.1f%%)\n", successful, successRate)
	fmt.Printf("❌ Failed:             %d (%.1f%%)\n", failed, 100-successRate)
	fmt.Printf("⏱️  Total Duration:     %v\n", duration)
	fmt.Printf("🚄 Requests/sec:       %.1f\n", rps)
	fmt.Printf("📊 Avg Latency:        %.1f ms\n", avgLatency)
	fmt.Printf("⚡ Min Latency:        %.1f ms\n", float64(hlt.stats.minLatency)/1000.0)
	fmt.Printf("🐌 Max Latency:        %.1f ms\n", float64(hlt.stats.maxLatency)/1000.0)

	fmt.Println("\n📋 Response Codes:")
	for code, count := range hlt.stats.responseCodes {
		fmt.Printf("   %d: %d requests\n", code, count)
	}
	fmt.Println("=")
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go <total_requests> [concurrency] [users] [duration_seconds] [server_url]")
		log.Println("Examples:")
		log.Println("  go run main.go 10000                    # 10k requests, default settings")
		log.Println("  go run main.go 50000 20                 # 50k requests, 20 concurrent")
		log.Println("  go run main.go 100000 50 5000 60        # 100k requests, 50 concurrent, 5k users, 60s max")
		log.Println("  go run main.go 10000 10 1000 30 http://localhost:8080  # Custom server")
		return
	}

	totalRequests := 10000
	if len(os.Args) > 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil {
			totalRequests = n
		}
	}

	concurrency := 10
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil {
			concurrency = n
		}
	}

	numUsers := 1000
	if len(os.Args) > 3 {
		if n, err := strconv.Atoi(os.Args[3]); err == nil {
			numUsers = n
		}
	}

	durationSeconds := 300 // 5 minutes default
	if len(os.Args) > 4 {
		if n, err := strconv.Atoi(os.Args[4]); err == nil {
			durationSeconds = n
		}
	}

	serverURL := "http://localhost:8080"
	if len(os.Args) > 5 {
		serverURL = os.Args[5]
	}

	duration := time.Duration(durationSeconds) * time.Second

	hlt := NewHTTPLoadTester(serverURL)
	if err := hlt.RunConcurrentLoadTest(totalRequests, concurrency, numUsers, duration); err != nil {
		log.Fatal("Load test failed:", err)
	}
}
