package botdetector

import "testing"

func TestIsBot(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		expected  bool
	}{
		// Search engine bots
		{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", true},
		{"Bingbot", "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", true},
		{"DuckDuckBot", "Mozilla/5.0 (compatible; DuckDuckBot-Https/1.1; https://duckduckgo.com/duckduckbot)", true},

		// Social media bots
		{"FacebookBot", "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)", true},
		{"TwitterBot", "Twitterbot/1.0", true},
		{"LinkedInBot", "LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)", true},

		// SEO bots
		{"AhrefsBot", "Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)", true},
		{"SemrushBot", "Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)", true},

		// HTTP libraries
		{"cURL", "curl/7.84.0", true},
		{"Wget", "Wget/1.21.3", true},
		{"Python Requests", "python-requests/2.28.1", true},
		{"Go HTTP Client", "Go-http-client/1.1", true},

		// Monitoring
		{"UptimeRobot", "Mozilla/5.0 (compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)", true},
		{"Pingdom", "Mozilla/5.0 (compatible; PingdomBot/1.0; +http://www.pingdom.com/)", true},

		// Real browsers (should be false)
		{"Chrome", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", false},
		{"Firefox", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0", false},
		{"Safari", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15", false},
		{"Edge", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0", false},

		// Mobile browsers
		{"iPhone Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1", false},
		{"Android Chrome", "Mozilla/5.0 (Linux; Android 14; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36", false},

		// Edge cases
		{"Empty user agent", "", true},
		{"WhatsApp", "WhatsApp/2.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBot(tt.userAgent)
			if result != tt.expected {
				t.Errorf("IsBot(%q) = %v, expected %v", tt.userAgent, result, tt.expected)
			}
		})
	}
}

func TestGetBotName(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		expected  string
	}{
		{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "Googlebot"},
		{"Bingbot", "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", "Bingbot"},
		{"FacebookBot", "facebookexternalhit/1.1", "Facebook Bot"},
		{"TwitterBot", "Twitterbot/1.0", "Twitter Bot"},
		{"AhrefsBot", "Mozilla/5.0 (compatible; AhrefsBot/7.0)", "Ahrefs Bot"},
		{"cURL", "curl/7.84.0", "cURL"},
		{"Python Requests", "python-requests/2.28.1", "Python Requests"},
		{"Generic Bot", "Some Random Bot", "Generic Bot"},
		{"Crawler", "My Crawler 1.0", "Crawler"},
		{"Empty", "", "Unknown Bot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBotName(tt.userAgent)
			if result != tt.expected {
				t.Errorf("GetBotName(%q) = %q, expected %q", tt.userAgent, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsBot(b *testing.B) {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"curl/7.84.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsBot(userAgents[i%len(userAgents)])
	}
}
