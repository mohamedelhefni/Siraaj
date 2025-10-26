package main

// High-Performance HTTP Load Tester for Analytics Service
//
// Optimizations:
// - Pre-generates all events to eliminate generation overhead during test
// - Uses large connection pool (1000 idle connections)
// - Disabled compression for lower latency
// - Increased buffer sizes (32KB read/write)
// - Removed artificial rate limiting - workers run at maximum speed
// - Minimized lock contention in stats collection
// - HTTP/1.1 for lower latency vs HTTP/2
// - Large event channel buffer (concurrency * 10)
//
// Expected Performance:
// - With Appender API: 10,000+ req/s on modern hardware
// - With regular INSERT: 1,000-2,000 req/s
//
// Usage:
//   go run http.go 100000 200        # 100k requests, 200 workers
//   go run http.go 500000 500        # 500k requests, 500 workers (stress test)

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

	// Bot user agents for simulation (about 20% bot traffic)
	botUserAgents = []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)",
		"Mozilla/5.0 (compatible; DuckDuckBot-Https/1.1; https://duckduckgo.com/duckduckbot)",
		"Twitterbot/1.0",
		"LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)",
		"curl/7.84.0",
		"python-requests/2.28.1",
		"Go-http-client/1.1",
		"Mozilla/5.0 (compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)",
		"Mozilla/5.0 (compatible; PingdomBot/1.0; +http://www.pingdom.com/)",
		"Wget/1.21.3",
		"WhatsApp/2.0",
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

func NewHTTPLoadTester(baseURL string) *HTTPLoadTester {
	return &HTTPLoadTester{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Reduced timeout for faster failure detection
			Transport: &http.Transport{
				MaxIdleConns:        1000, // Increased for high concurrency
				MaxIdleConnsPerHost: 1000, // Match MaxIdleConns
				MaxConnsPerHost:     0,    // No limit
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,     // Keep connections alive
				DisableCompression:  true,      // Disable compression for speed
				WriteBufferSize:     32 * 1024, // 32KB write buffer
				ReadBufferSize:      32 * 1024, // 32KB read buffer
				ForceAttemptHTTP2:   false,     // Use HTTP/1.1 for lower latency
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

	// 20% chance of bot traffic
	var userAgent string
	if rand.Float32() < 0.20 {
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
		Timestamp: timestamp,
		EventName: eventName,
		UserID:    userID,
		SessionID: sessionID,
		URL:       url,
		Referrer:  referrer,
		UserAgent: userAgent,
		IP:        ip,
		Country:   country,
		Browser:   browser,
		OS:        os,
		Device:    device,
	}
}

// generateProperties creates realistic properties based on event type (deprecated - no longer used)
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

// updateStats safely updates the performance statistics with minimal locking
func (hlt *HTTPLoadTester) updateStats(latency int64, resp *http.Response, err error) {
	// Use atomics for counters (no lock needed)
	atomic.AddInt64(&hlt.stats.totalRequests, 1)
	atomic.AddInt64(&hlt.stats.totalLatency, latency)

	if err != nil {
		atomic.AddInt64(&hlt.stats.failedReqs, 1)
	} else {
		atomic.AddInt64(&hlt.stats.successfulReqs, 1)
	}

	// Only lock for response codes and min/max latency
	hlt.stats.mutex.Lock()

	if resp != nil {
		hlt.stats.responseCodes[resp.StatusCode]++
	}

	// Update min/max latency
	if latency < hlt.stats.minLatency {
		hlt.stats.minLatency = latency
	}
	if latency > hlt.stats.maxLatency {
		hlt.stats.maxLatency = latency
	}

	hlt.stats.mutex.Unlock()
}

// RunConcurrentLoadTest runs a high-performance concurrent HTTP load test
func (hlt *HTTPLoadTester) RunConcurrentLoadTest(totalRequests int, concurrency int, numUsers int, duration time.Duration) error {
	log.Printf("🚀 Starting HIGH-PERFORMANCE HTTP load test")
	log.Printf("📊 Target: %s/api/track", hlt.baseURL)
	log.Printf("📈 Total requests: %d", totalRequests)
	log.Printf("🔄 Concurrency: %d workers", concurrency)
	log.Printf("👥 Users: %d", numUsers)
	log.Printf("⏱️  Max Duration: %v", duration)

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

	// Pre-generate all events for maximum throughput
	log.Printf("🔄 Pre-generating %d events...", totalRequests)
	events := make([]Event, totalRequests)
	for i := 0; i < totalRequests; i++ {
		events[i] = hlt.GenerateRandomEvent(userPool)
	}
	log.Printf("✅ Events generated, starting load test...")

	hlt.stats.startTime = time.Now()

	// Large buffer channel for maximum throughput
	eventChan := make(chan Event, concurrency*10)
	var wg sync.WaitGroup

	// Atomic counter for completed requests
	var completed int64

	// Start worker goroutines - each worker processes events as fast as possible
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for event := range eventChan {
				if err := hlt.SendEvent(event); err != nil {
					// Don't log every error, just count them
					if atomic.LoadInt64(&hlt.stats.failedReqs)%100 == 0 {
						log.Printf("Worker %d: %d errors so far", workerID, atomic.LoadInt64(&hlt.stats.failedReqs))
					}
				}
				atomic.AddInt64(&completed, 1)
			}
		}(i)
	}

	// Feed events to workers as fast as possible
	go func() {
		defer close(eventChan)

		timeout := time.After(duration)
		sent := 0

		for sent < totalRequests {
			select {
			case <-timeout:
				log.Printf("⏰ Duration limit reached after %d requests", sent)
				return
			case eventChan <- events[sent]:
				sent++
				// No artificial rate limiting - send as fast as workers can handle
			}
		}
		log.Printf("🎯 All %d requests queued", sent)
	}()

	// Progress reporting
	progressDone := make(chan struct{})
	go func() {
		progressTicker := time.NewTicker(2 * time.Second)
		defer progressTicker.Stop()

		lastCount := int64(0)
		lastTime := time.Now()

		for {
			select {
			case <-progressTicker.C:
				currentCount := atomic.LoadInt64(&completed)
				currentTime := time.Now()

				// Calculate instantaneous RPS
				countDelta := currentCount - lastCount
				timeDelta := currentTime.Sub(lastTime).Seconds()
				instantRPS := float64(countDelta) / timeDelta

				hlt.printProgressWithRPS(instantRPS)

				lastCount = currentCount
				lastTime = currentTime
			case <-progressDone:
				return
			}
		}
	}()

	// Wait for all workers to complete
	wg.Wait()
	close(progressDone)
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

// printProgressWithRPS displays progress with instantaneous RPS
func (hlt *HTTPLoadTester) printProgressWithRPS(instantRPS float64) {
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
	overallRPS := float64(total) / elapsed.Seconds()

	log.Printf("📊 %d reqs | ✅ %d | ❌ %d | 🚄 %.0f req/s (inst: %.0f) | ⏱️ %.1fms avg",
		total, successful, failed, overallRPS, instantRPS, float64(avgLatency)/1000.0)
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
		log.Println("Usage: go run http.go <total_requests> [concurrency] [users] [duration_seconds] [server_url]")
		log.Println("Examples:")
		log.Println("  go run http.go 10000                    # 10k requests, default settings")
		log.Println("  go run http.go 100000 100               # 100k requests, 100 workers")
		log.Println("  go run http.go 500000 200 5000 300      # 500k requests, 200 workers, 5k users, 300s")
		log.Println("  go run http.go 10000 50 1000 30 http://localhost:8080  # Custom server")
		log.Println("\nRecommended for high performance:")
		log.Println("  go run http.go 100000 200               # 200 concurrent workers")
		return
	}

	totalRequests := 100000 // Increased default
	if len(os.Args) > 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil {
			totalRequests = n
		}
	}

	concurrency := 100 // Increased default for better throughput
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil {
			concurrency = n
		}
	}

	numUsers := 10000 // Increased default
	if len(os.Args) > 3 {
		if n, err := strconv.Atoi(os.Args[3]); err == nil {
			numUsers = n
		}
	}

	durationSeconds := 600 // 10 minutes default (increased)
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

	log.Printf("⚙️  Configuration:")
	log.Printf("   Requests:    %d", totalRequests)
	log.Printf("   Concurrency: %d", concurrency)
	log.Printf("   Users:       %d", numUsers)
	log.Printf("   Duration:    %v", duration)
	log.Printf("   Server:      %s\n", serverURL)

	hlt := NewHTTPLoadTester(serverURL)
	if err := hlt.RunConcurrentLoadTest(totalRequests, concurrency, numUsers, duration); err != nil {
		log.Fatal("Load test failed:", err)
	}
}
