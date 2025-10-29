package channeldetector

import (
	"net/url"
	"strings"
)

// Channel represents the traffic source category
type Channel string

const (
	ChannelDirect   Channel = "Direct"
	ChannelOrganic  Channel = "Organic"
	ChannelReferral Channel = "Referral"
	ChannelSocial   Channel = "Social"
	ChannelPaid     Channel = "Paid"
)

// Common organic search engines
var organicSearchEngines = []string{
	"google.com",
	"bing.com",
	"yahoo.com",
	"duckduckgo.com",
	"baidu.com",
	"yandex.com",
	"ask.com",
	"aol.com",
	"ecosia.org",
}

// Common social media platforms
var socialPlatforms = []string{
	"facebook.com",
	"twitter.com",
	"linkedin.com",
	"instagram.com",
	"pinterest.com",
	"reddit.com",
	"tiktok.com",
	"youtube.com",
	"snapchat.com",
	"tumblr.com",
	"vk.com",
	"weibo.com",
	"t.co",  // Twitter shortener
	"fb.me", // Facebook shortener
}

// Common paid ad parameters
var paidParameters = []string{
	"utm_source=google",
	"utm_source=facebook",
	"utm_source=bing",
	"utm_source=linkedin",
	"utm_source=twitter",
	"utm_medium=cpc",
	"utm_medium=ppc",
	"utm_medium=paid",
	"utm_medium=display",
	"gclid=",     // Google Ads click ID
	"fbclid=",    // Facebook click ID
	"msclkid=",   // Microsoft Ads click ID
	"twclid=",    // Twitter click ID
	"li_fat_id=", // LinkedIn click ID
}

// DetectChannel classifies an event based on its referrer and URL
// Returns one of: Direct, Organic, Referral, Social, Paid
func DetectChannel(referrer string, eventURL string, currentDomain string) Channel {
	// Clean up inputs
	referrer = strings.TrimSpace(strings.ToLower(referrer))
	eventURL = strings.TrimSpace(strings.ToLower(eventURL))
	currentDomain = strings.TrimSpace(strings.ToLower(currentDomain))

	// Check for paid traffic first (highest priority)
	if isPaid(eventURL) {
		return ChannelPaid
	}

	// Check if referrer is empty (Direct traffic)
	if referrer == "" || referrer == "(direct)" || referrer == "null" {
		return ChannelDirect
	}

	// Parse referrer to get domain
	referrerDomain := extractDomain(referrer)
	if referrerDomain == "" {
		return ChannelDirect
	}

	// Check if referrer is the same as current domain (internal navigation = Direct)
	if currentDomain != "" && strings.Contains(referrerDomain, currentDomain) {
		return ChannelDirect
	}

	// Check for social media
	if isSocial(referrerDomain) {
		return ChannelSocial
	}

	// Check for organic search
	if isOrganic(referrerDomain) {
		return ChannelOrganic
	}

	// Everything else is referral traffic
	return ChannelReferral
}

// isPaid checks if the URL contains paid advertising parameters
func isPaid(eventURL string) bool {
	lowerURL := strings.ToLower(eventURL)
	for _, param := range paidParameters {
		if strings.Contains(lowerURL, param) {
			return true
		}
	}
	return false
}

// isSocial checks if a domain is a social media platform
func isSocial(domain string) bool {
	for _, social := range socialPlatforms {
		if strings.Contains(domain, social) {
			return true
		}
	}
	return false
}

// isOrganic checks if a domain is an organic search engine
func isOrganic(domain string) bool {
	for _, engine := range organicSearchEngines {
		if strings.Contains(domain, engine) {
			return true
		}
	}
	return false
}

// extractDomain extracts the domain from a URL string
func extractDomain(urlStr string) string {
	// Handle cases where URL doesn't have a scheme
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	return parsedURL.Hostname()
}
