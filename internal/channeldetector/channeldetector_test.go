package channeldetector

import (
	"testing"
)

func TestDetectChannel(t *testing.T) {
	tests := []struct {
		name          string
		referrer      string
		eventURL      string
		currentDomain string
		expected      Channel
	}{
		{
			name:          "Direct - empty referrer",
			referrer:      "",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelDirect,
		},
		{
			name:          "Direct - null referrer",
			referrer:      "null",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelDirect,
		},
		{
			name:          "Direct - same domain",
			referrer:      "https://example.com/other-page",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelDirect,
		},
		{
			name:          "Organic - Google",
			referrer:      "https://www.google.com/search",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelOrganic,
		},
		{
			name:          "Organic - Bing",
			referrer:      "https://www.bing.com/search?q=test",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelOrganic,
		},
		{
			name:          "Organic - Yahoo",
			referrer:      "https://search.yahoo.com/search",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelOrganic,
		},
		{
			name:          "Social - Facebook",
			referrer:      "https://www.facebook.com",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelSocial,
		},
		{
			name:          "Social - Twitter",
			referrer:      "https://t.co/abc123",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelSocial,
		},
		{
			name:          "Social - LinkedIn",
			referrer:      "https://www.linkedin.com/feed",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelSocial,
		},
		{
			name:          "Social - Instagram",
			referrer:      "https://www.instagram.com",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelSocial,
		},
		{
			name:          "Paid - Google Ads (gclid)",
			referrer:      "https://www.google.com",
			eventURL:      "https://example.com/page?gclid=abc123",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
		{
			name:          "Paid - Facebook Ads (fbclid)",
			referrer:      "https://www.facebook.com",
			eventURL:      "https://example.com/page?fbclid=xyz789",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
		{
			name:          "Paid - UTM medium cpc",
			referrer:      "https://newsletter.example.org",
			eventURL:      "https://example.com/page?utm_source=google&utm_medium=cpc",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
		{
			name:          "Paid - UTM medium ppc",
			referrer:      "",
			eventURL:      "https://example.com/page?utm_medium=ppc&utm_campaign=summer",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
		{
			name:          "Referral - other website",
			referrer:      "https://news.ycombinator.com",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelReferral,
		},
		{
			name:          "Referral - blog",
			referrer:      "https://blog.somesite.com/article",
			eventURL:      "https://example.com/page",
			currentDomain: "example.com",
			expected:      ChannelReferral,
		},
		{
			name:          "Paid takes precedence over Social",
			referrer:      "https://www.facebook.com",
			eventURL:      "https://example.com/page?utm_medium=cpc",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
		{
			name:          "Paid takes precedence over Organic",
			referrer:      "https://www.google.com",
			eventURL:      "https://example.com/page?gclid=abc123",
			currentDomain: "example.com",
			expected:      ChannelPaid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectChannel(tt.referrer, tt.eventURL, tt.currentDomain)
			if result != tt.expected {
				t.Errorf("DetectChannel(%q, %q, %q) = %v, want %v",
					tt.referrer, tt.eventURL, tt.currentDomain, result, tt.expected)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Full URL with https",
			url:      "https://www.example.com/path",
			expected: "www.example.com",
		},
		{
			name:     "Full URL with http",
			url:      "http://example.com",
			expected: "example.com",
		},
		{
			name:     "URL without scheme",
			url:      "www.google.com/search",
			expected: "www.google.com",
		},
		{
			name:     "Subdomain",
			url:      "https://blog.example.com",
			expected: "blog.example.com",
		},
		{
			name:     "URL with port",
			url:      "https://example.com:8080/path",
			expected: "example.com",
		},
		{
			name:     "URL with query params",
			url:      "https://example.com/page?param=value",
			expected: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDomain(tt.url)
			if result != tt.expected {
				t.Errorf("extractDomain(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestIsPaid(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"Google Ads click ID", "https://example.com?gclid=abc123", true},
		{"Facebook click ID", "https://example.com?fbclid=xyz789", true},
		{"UTM medium CPC", "https://example.com?utm_medium=cpc", true},
		{"UTM medium PPC", "https://example.com?utm_medium=ppc", true},
		{"UTM source Google", "https://example.com?utm_source=google&utm_medium=email", true},
		{"No paid params", "https://example.com", false},
		{"Organic params", "https://example.com?ref=blog", false},
		{"Case insensitive", "https://example.com?GCLID=ABC123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPaid(tt.url)
			if result != tt.expected {
				t.Errorf("isPaid(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestIsSocial(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{"Facebook", "facebook.com", true},
		{"Facebook www", "www.facebook.com", true},
		{"Twitter", "twitter.com", true},
		{"Twitter shortener", "t.co", true},
		{"LinkedIn", "linkedin.com", true},
		{"Instagram", "instagram.com", true},
		{"Reddit", "reddit.com", true},
		{"YouTube", "youtube.com", true},
		{"Not social", "example.com", false},
		{"Google", "google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSocial(tt.domain)
			if result != tt.expected {
				t.Errorf("isSocial(%q) = %v, want %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestIsOrganic(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{"Google", "google.com", true},
		{"Google www", "www.google.com", true},
		{"Bing", "bing.com", true},
		{"Yahoo", "yahoo.com", true},
		{"DuckDuckGo", "duckduckgo.com", true},
		{"Facebook", "facebook.com", false},
		{"Random site", "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOrganic(tt.domain)
			if result != tt.expected {
				t.Errorf("isOrganic(%q) = %v, want %v", tt.domain, result, tt.expected)
			}
		})
	}
}
