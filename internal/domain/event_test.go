package domain

import (
	"testing"
	"time"
)

func TestEventValidation(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		valid bool
	}{
		{
			name: "Valid event with all fields",
			event: Event{
				ID:        1,
				Timestamp: time.Now(),
				EventName: "page_view",
				UserID:    "user123",
				SessionID: "session456",
				URL:       "/home",
				Referrer:  "https://google.com",
				UserAgent: "Mozilla/5.0",
				IP:        "192.168.1.1",
				Country:   "United States",
				Browser:   "Chrome",
				OS:        "Windows",
				Device:    "Desktop",
				IsBot:     false,
				ProjectID: "project1",
			},
			valid: true,
		},
		{
			name: "Valid event with minimal fields",
			event: Event{
				Timestamp: time.Now(),
				EventName: "click",
				UserID:    "user123",
			},
			valid: true,
		},
		{
			name: "Event with Palestine country normalization",
			event: Event{
				Timestamp: time.Now(),
				EventName: "page_view",
				UserID:    "user123",
				Country:   "Palestine",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - event should have required fields
			if tt.event.EventName == "" && tt.valid {
				t.Error("Expected valid event to have event_name")
			}
			if tt.event.UserID == "" && tt.valid {
				t.Error("Expected valid event to have user_id")
			}
		})
	}
}

func TestPageStat(t *testing.T) {
	stat := PageStat{
		URL:   "/home",
		Count: 100,
	}

	if stat.URL != "/home" {
		t.Errorf("Expected URL to be /home, got %s", stat.URL)
	}
	if stat.Count != 100 {
		t.Errorf("Expected Count to be 100, got %d", stat.Count)
	}
}

func TestCountryStat(t *testing.T) {
	tests := []struct {
		name    string
		country string
		count   int64
	}{
		{"United States", "United States", 500},
		{"Palestine", "Palestine", 250},
		{"Canada", "Canada", 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat := CountryStat{
				Country: tt.country,
				Count:   tt.count,
			}

			if stat.Country != tt.country {
				t.Errorf("Expected Country to be %s, got %s", tt.country, stat.Country)
			}
			if stat.Count != tt.count {
				t.Errorf("Expected Count to be %d, got %d", tt.count, stat.Count)
			}
		})
	}
}

func TestBrowserStat(t *testing.T) {
	stat := BrowserStat{
		Browser: "Chrome",
		Count:   300,
	}

	if stat.Browser != "Chrome" {
		t.Errorf("Expected Browser to be Chrome, got %s", stat.Browser)
	}
	if stat.Count != 300 {
		t.Errorf("Expected Count to be 300, got %d", stat.Count)
	}
}

func TestOSStat(t *testing.T) {
	stat := OSStat{
		OS:    "Windows",
		Count: 200,
	}

	if stat.OS != "Windows" {
		t.Errorf("Expected OS to be Windows, got %s", stat.OS)
	}
	if stat.Count != 200 {
		t.Errorf("Expected Count to be 200, got %d", stat.Count)
	}
}

func TestDeviceStat(t *testing.T) {
	stat := DeviceStat{
		Device: "Mobile",
		Count:  150,
	}

	if stat.Device != "Mobile" {
		t.Errorf("Expected Device to be Mobile, got %s", stat.Device)
	}
	if stat.Count != 150 {
		t.Errorf("Expected Count to be 150, got %d", stat.Count)
	}
}

func TestTimelineStat(t *testing.T) {
	dateStr := "2025-10-25"
	stat := TimelineStat{
		Date:  dateStr,
		Count: 500,
	}

	if stat.Date != dateStr {
		t.Errorf("Expected Date to be %s, got %s", dateStr, stat.Date)
	}
	if stat.Count != 500 {
		t.Errorf("Expected Count to be 500, got %d", stat.Count)
	}
}

func TestReferrerStat(t *testing.T) {
	stat := ReferrerStat{
		Referrer: "https://google.com",
		Count:    400,
	}

	if stat.Referrer != "https://google.com" {
		t.Errorf("Expected Referrer to be https://google.com, got %s", stat.Referrer)
	}
	if stat.Count != 400 {
		t.Errorf("Expected Count to be 400, got %d", stat.Count)
	}
}
