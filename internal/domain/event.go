package domain

import "time"

type Event struct {
	ID         uint64    `json:"id"`
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
	IsBot      bool      `json:"is_bot"`
	Properties string    `json:"properties"` // JSON string
	ProjectID  string    `json:"project_id"`
}

type Stats struct {
	PageViews      int64            `json:"page_views"`
	UniqueVisitors int64            `json:"unique_visitors"`
	UniqueUsers    int64            `json:"unique_users"`
	TopPages       []PageStat       `json:"top_pages"`
	TopReferrers   []ReferrerStat   `json:"top_referrers"`
	Countries      []CountryStat    `json:"countries"`
	Browsers       []BrowserStat    `json:"browsers"`
	OSList         []OSStat         `json:"os"`
	Devices        []DeviceStat     `json:"devices"`
	Timeline       []TimelineStat   `json:"timeline"`
	Events         map[string]int64 `json:"events"`
}

type PageStat struct {
	URL   string `json:"url"`
	Count int64  `json:"count"`
}

type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type CountryStat struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type BrowserStat struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

type OSStat struct {
	OS    string `json:"os"`
	Count int64  `json:"count"`
}

type DeviceStat struct {
	Device string `json:"device"`
	Count  int64  `json:"count"`
}

type TimelineStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type OnlineUsers struct {
	Count int64 `json:"count"`
}

type Project struct {
	ID         string `json:"id"`
	EventCount int64  `json:"event_count"`
}
