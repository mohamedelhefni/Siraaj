package botdetector

import (
	"regexp"
	"strings"
)

// Common bot user agent patterns
var botPatterns = []*regexp.Regexp{
	// Search engine bots
	regexp.MustCompile(`(?i)googlebot`),
	regexp.MustCompile(`(?i)bingbot`),
	regexp.MustCompile(`(?i)yahoo`),
	regexp.MustCompile(`(?i)duckduckbot`),
	regexp.MustCompile(`(?i)baiduspider`),
	regexp.MustCompile(`(?i)yandex`),
	regexp.MustCompile(`(?i)slurp`), // Yahoo Slurp

	// Social media bots
	regexp.MustCompile(`(?i)facebookexternalhit`),
	regexp.MustCompile(`(?i)twitterbot`),
	regexp.MustCompile(`(?i)linkedinbot`),
	regexp.MustCompile(`(?i)whatsapp`),
	regexp.MustCompile(`(?i)telegrambot`),
	regexp.MustCompile(`(?i)discordbot`),
	regexp.MustCompile(`(?i)slackbot`),

	// SEO/Monitoring bots
	regexp.MustCompile(`(?i)ahrefsbot`),
	regexp.MustCompile(`(?i)semrushbot`),
	regexp.MustCompile(`(?i)mj12bot`), // Majestic
	regexp.MustCompile(`(?i)dotbot`),
	regexp.MustCompile(`(?i)rogerbot`),
	regexp.MustCompile(`(?i)screaming frog`),
	regexp.MustCompile(`(?i)sitebulb`),

	// Generic bot indicators
	regexp.MustCompile(`(?i)\bbot\b`),
	regexp.MustCompile(`(?i)\bcrawler\b`),
	regexp.MustCompile(`(?i)\bspider\b`),
	regexp.MustCompile(`(?i)\bscraper\b`),
	regexp.MustCompile(`(?i)\bfetcher\b`),

	// Headless browsers (often used for scraping)
	regexp.MustCompile(`(?i)headlesschrome`),
	regexp.MustCompile(`(?i)phantomjs`),
	regexp.MustCompile(`(?i)selenium`),
	regexp.MustCompile(`(?i)webdriver`),
	regexp.MustCompile(`(?i)puppeteer`),

	// Monitoring services
	regexp.MustCompile(`(?i)pingdom`),
	regexp.MustCompile(`(?i)uptimerobot`),
	regexp.MustCompile(`(?i)newrelic`),
	regexp.MustCompile(`(?i)statuscake`),
	regexp.MustCompile(`(?i)sitechecker`),

	// Archiving/Indexing
	regexp.MustCompile(`(?i)archive\.org`),
	regexp.MustCompile(`(?i)ia_archiver`),
	regexp.MustCompile(`(?i)wayback`),

	// HTTP libraries
	regexp.MustCompile(`(?i)^curl`),
	regexp.MustCompile(`(?i)^wget`),
	regexp.MustCompile(`(?i)^python-requests`),
	regexp.MustCompile(`(?i)^go-http-client`),
	regexp.MustCompile(`(?i)^axios`),
	regexp.MustCompile(`(?i)^httpie`),
}

// Additional suspicious patterns
var suspiciousPatterns = []string{
	"http",    // HTTP libraries often don't include full user agents
	"library", // Generic library indicators
	"fetcher", // Data fetching tools
	"monitoring",
	"check",
}

// IsBot determines if a user agent string belongs to a bot
func IsBot(userAgent string) bool {
	if userAgent == "" {
		return true // Empty user agent is suspicious
	}

	// Normalize user agent
	ua := strings.TrimSpace(userAgent)

	// Check against known bot patterns
	for _, pattern := range botPatterns {
		if pattern.MatchString(ua) {
			return true
		}
	}

	// Additional heuristics for suspicious user agents
	uaLower := strings.ToLower(ua)

	// Very short user agents are often bots
	if len(ua) < 20 {
		for _, suspicious := range suspiciousPatterns {
			if strings.Contains(uaLower, suspicious) {
				return true
			}
		}
	}

	// Check for missing common browser indicators
	hasCommonBrowser := strings.Contains(uaLower, "mozilla") ||
		strings.Contains(uaLower, "chrome") ||
		strings.Contains(uaLower, "safari") ||
		strings.Contains(uaLower, "firefox") ||
		strings.Contains(uaLower, "edge")

	// If it doesn't look like a browser and contains suspicious keywords
	if !hasCommonBrowser {
		for _, suspicious := range suspiciousPatterns {
			if strings.Contains(uaLower, suspicious) {
				return true
			}
		}
	}

	return false
}

// GetBotName attempts to identify the specific bot name
func GetBotName(userAgent string) string {
	if userAgent == "" {
		return "Unknown Bot"
	}

	uaLower := strings.ToLower(userAgent)

	// Search engine bots
	if strings.Contains(uaLower, "googlebot") {
		return "Googlebot"
	}
	if strings.Contains(uaLower, "bingbot") {
		return "Bingbot"
	}
	if strings.Contains(uaLower, "duckduckbot") {
		return "DuckDuckBot"
	}
	if strings.Contains(uaLower, "baiduspider") {
		return "Baidu Spider"
	}
	if strings.Contains(uaLower, "yandex") {
		return "Yandex Bot"
	}
	if strings.Contains(uaLower, "slurp") {
		return "Yahoo Slurp"
	}

	// Social media bots
	if strings.Contains(uaLower, "facebookexternalhit") {
		return "Facebook Bot"
	}
	if strings.Contains(uaLower, "twitterbot") {
		return "Twitter Bot"
	}
	if strings.Contains(uaLower, "linkedinbot") {
		return "LinkedIn Bot"
	}
	if strings.Contains(uaLower, "whatsapp") {
		return "WhatsApp Bot"
	}
	if strings.Contains(uaLower, "telegrambot") {
		return "Telegram Bot"
	}

	// SEO bots
	if strings.Contains(uaLower, "ahrefsbot") {
		return "Ahrefs Bot"
	}
	if strings.Contains(uaLower, "semrushbot") {
		return "SEMrush Bot"
	}
	if strings.Contains(uaLower, "mj12bot") {
		return "Majestic Bot"
	}

	// Monitoring
	if strings.Contains(uaLower, "pingdom") {
		return "Pingdom"
	}
	if strings.Contains(uaLower, "uptimerobot") {
		return "UptimeRobot"
	}

	// HTTP libraries
	if strings.Contains(uaLower, "curl") {
		return "cURL"
	}
	if strings.Contains(uaLower, "wget") {
		return "Wget"
	}
	if strings.Contains(uaLower, "python-requests") {
		return "Python Requests"
	}
	if strings.Contains(uaLower, "go-http-client") {
		return "Go HTTP Client"
	}

	// Generic patterns
	if strings.Contains(uaLower, "bot") {
		return "Generic Bot"
	}
	if strings.Contains(uaLower, "crawler") {
		return "Crawler"
	}
	if strings.Contains(uaLower, "spider") {
		return "Spider"
	}

	return "Unknown Bot"
}
