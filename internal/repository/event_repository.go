package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mohamedelhefni/siraaj/internal/domain"
)

type EventRepository interface {
	Create(event domain.Event) error
	CreateBatch(events []domain.Event) error
	GetEvents(startDate, endDate time.Time, limit, offset int) (map[string]interface{}, error)
	GetStats(startDate, endDate time.Time, limit int, filters map[string]string) (map[string]interface{}, error)
	GetOnlineUsers(timeWindow int) (map[string]interface{}, error)
	GetProjects() ([]string, error)
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event domain.Event) error {
	query := `
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if event.ProjectID == "" {
		event.ProjectID = "default"
	}

	_, err := r.db.Exec(query,
		event.Timestamp,
		event.EventName,
		event.UserID,
		event.SessionID,
		event.URL,
		event.Referrer,
		event.UserAgent,
		event.IP,
		event.Country,
		event.Browser,
		event.OS,
		event.Device,
		event.Properties,
		event.ProjectID,
	)

	if err != nil {
		log.Printf("Error inserting event: %v", err)
	}

	return err
}

func (r *eventRepository) CreateBatch(events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, url, referrer, 
			user_agent, ip, country, browser, os, device, properties, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		if event.ProjectID == "" {
			event.ProjectID = "default"
		}

		_, err := stmt.Exec(
			event.Timestamp,
			event.EventName,
			event.UserID,
			event.SessionID,
			event.URL,
			event.Referrer,
			event.UserAgent,
			event.IP,
			event.Country,
			event.Browser,
			event.OS,
			event.Device,
			event.Properties,
			event.ProjectID,
		)
		if err != nil {
			log.Printf("Error inserting event in batch: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (r *eventRepository) GetEvents(startDate, endDate time.Time, limit, offset int) (map[string]interface{}, error) {
	query := `
		SELECT id, timestamp, event_name, user_id, session_id, url, referrer,
			user_agent, ip, country, browser, os, device, properties, project_id
		FROM events
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		err := rows.Scan(
			&e.ID, &e.Timestamp, &e.EventName, &e.UserID, &e.SessionID,
			&e.URL, &e.Referrer, &e.UserAgent, &e.IP, &e.Country,
			&e.Browser, &e.OS, &e.Device, &e.Properties, &e.ProjectID,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			continue
		}
		events = append(events, e)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM events WHERE timestamp BETWEEN ? AND ?`
	err = r.db.QueryRow(countQuery, startDate, endDate).Scan(&total)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, nil
}

func (r *eventRepository) GetStats(startDate, endDate time.Time, limit int, filters map[string]string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	if limit <= 0 {
		limit = 10
	}

	// Build WHERE clause based on filters
	whereClause := "timestamp BETWEEN ? AND ?"
	args := []interface{}{startDate, endDate}

	if projectID, ok := filters["project"]; ok && projectID != "" {
		whereClause += " AND project_id = ?"
		args = append(args, projectID)
	}
	if source, ok := filters["source"]; ok && source != "" {
		whereClause += " AND referrer = ?"
		args = append(args, source)
	}
	if country, ok := filters["country"]; ok && country != "" {
		whereClause += " AND country = ?"
		args = append(args, country)
	}
	if browser, ok := filters["browser"]; ok && browser != "" {
		whereClause += " AND browser = ?"
		args = append(args, browser)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		whereClause += " AND event_name = ?"
		args = append(args, eventName)
	}

	// Use a single query with CTEs for better performance
	aggregateQuery := fmt.Sprintf(`
	WITH date_filtered AS (
		SELECT * FROM events 
		WHERE %s
	),
	event_stats AS (
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT session_id) as total_visits,
			COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_views,
			COUNT(DISTINCT CASE WHEN event_name = 'page_view' THEN session_id END) as sessions_with_views
		FROM date_filtered
	)
	SELECT total_events, unique_users, total_visits, page_views, sessions_with_views FROM event_stats;
	`, whereClause)

	var totalEvents, uniqueUsers, totalVisits, pageViews, sessionsWithViews int
	err := r.db.QueryRow(aggregateQuery, args...).Scan(&totalEvents, &uniqueUsers, &totalVisits, &pageViews, &sessionsWithViews)
	if err != nil {
		return nil, err
	}
	stats["total_events"] = totalEvents
	stats["unique_users"] = uniqueUsers
	stats["total_visits"] = totalVisits
	stats["page_views"] = pageViews

	// Calculate bounce rate: sessions with only 1 page view / total sessions
	var bounceRate float64
	if totalVisits > 0 {
		singlePageQuery := fmt.Sprintf(`
			SELECT COUNT(DISTINCT session_id) 
			FROM (
				SELECT session_id, COUNT(*) as view_count
				FROM events 
				WHERE %s AND event_name = 'page_view'
				GROUP BY session_id
				HAVING view_count = 1
			)
		`, whereClause)
		var singlePageSessions int
		err = r.db.QueryRow(singlePageQuery, args...).Scan(&singlePageSessions)
		if err == nil && sessionsWithViews > 0 {
			bounceRate = float64(singlePageSessions) / float64(sessionsWithViews) * 100
		}
	}
	stats["bounce_rate"] = bounceRate

	// Top events with optimized query
	query := fmt.Sprintf(`
		SELECT event_name, COUNT(*) as count 
		FROM events 
		WHERE %s
		GROUP BY event_name 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)
	queryArgs := append(args, limit)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topEvents := []map[string]interface{}{}
	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			continue
		}
		topEvents = append(topEvents, map[string]interface{}{
			"name":  name,
			"count": count,
		})
	}
	stats["top_events"] = topEvents

	// Events over time with dynamic granularity based on date range
	timelineDuration := endDate.Sub(startDate)
	var timelineQuery string
	var timeFormat string

	// Determine granularity based on date range
	if timelineDuration <= 24*time.Hour {
		// For today or single day: show hourly data
		timelineQuery = fmt.Sprintf(`
			SELECT 
				strftime(DATE_TRUNC('hour', timestamp), '%%Y-%%m-%%d %%H:00:00') as date, 
				COUNT(*) as count
			FROM events 
			WHERE %s
			GROUP BY date 
			ORDER BY date
		`, whereClause)
		timeFormat = "hour"
	} else if timelineDuration <= 90*24*time.Hour {
		// For up to 3 months: show daily data
		timelineQuery = fmt.Sprintf(`
			SELECT 
				strftime(DATE_TRUNC('day', timestamp), '%%Y-%%m-%%d') as date, 
				COUNT(*) as count
			FROM events 
			WHERE %s
			GROUP BY date 
			ORDER BY date
		`, whereClause)
		timeFormat = "day"
	} else {
		// For more than 3 months: show monthly data
		timelineQuery = fmt.Sprintf(`
			SELECT 
				strftime(DATE_TRUNC('month', timestamp), '%%Y-%%m-01') as date, 
				COUNT(*) as count
			FROM events 
			WHERE %s
			GROUP BY date 
			ORDER BY date
		`, whereClause)
		timeFormat = "month"
	}

	rows, err = r.db.Query(timelineQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timeline := []map[string]interface{}{}
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			continue
		}
		timeline = append(timeline, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	stats["timeline"] = timeline
	stats["timeline_format"] = timeFormat

	// Top pages
	query = fmt.Sprintf(`
		SELECT url, COUNT(*) as count 
		FROM events 
		WHERE %s AND url IS NOT NULL AND url != ''
		GROUP BY url 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topPages := []map[string]interface{}{}
	for rows.Next() {
		var url string
		var count int
		if err := rows.Scan(&url, &count); err != nil {
			continue
		}
		topPages = append(topPages, map[string]interface{}{
			"url":   url,
			"count": count,
		})
	}
	stats["top_pages"] = topPages

	// Browsers
	query = fmt.Sprintf(`
		SELECT browser, COUNT(*) as count 
		FROM events 
		WHERE %s AND browser IS NOT NULL AND browser != ''
		GROUP BY browser 
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	rows, err = r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	browsers := []map[string]interface{}{}
	for rows.Next() {
		var browser string
		var count int
		if err := rows.Scan(&browser, &count); err != nil {
			continue
		}
		browsers = append(browsers, map[string]interface{}{
			"name":  browser,
			"count": count,
		})
	}
	stats["browsers"] = browsers

	// Top Countries
	query = fmt.Sprintf(`
		SELECT country, COUNT(*) as count 
		FROM events 
		WHERE %s AND country IS NOT NULL AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topCountries := []map[string]interface{}{}
	for rows.Next() {
		var country string
		var count int
		if err := rows.Scan(&country, &count); err != nil {
			continue
		}
		topCountries = append(topCountries, map[string]interface{}{
			"name":  country,
			"count": count,
		})
	}
	stats["top_countries"] = topCountries

	// Top Sources (Referrers) with URL parsing
	query = fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN referrer = '' OR referrer IS NULL THEN 'Direct'
				ELSE referrer
			END as source,
			COUNT(*) as count 
		FROM events 
		WHERE %s
		GROUP BY source 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	rows, err = r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topSources := []map[string]interface{}{}
	for rows.Next() {
		var referrer string
		var count int
		if err := rows.Scan(&referrer, &count); err != nil {
			continue
		}
		topSources = append(topSources, map[string]interface{}{
			"name":  referrer,
			"count": count,
		})
	}
	stats["top_sources"] = topSources

	// Calculate trends by comparing with previous period
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := startDate

	prevWhereClause := "timestamp BETWEEN ? AND ?"
	prevArgs := []interface{}{prevStartDate, prevEndDate}

	// Apply same filters to previous period
	if projectID, ok := filters["project"]; ok && projectID != "" {
		prevWhereClause += " AND project_id = ?"
		prevArgs = append(prevArgs, projectID)
	}
	if source, ok := filters["source"]; ok && source != "" {
		prevWhereClause += " AND referrer = ?"
		prevArgs = append(prevArgs, source)
	}
	if country, ok := filters["country"]; ok && country != "" {
		prevWhereClause += " AND country = ?"
		prevArgs = append(prevArgs, country)
	}
	if browser, ok := filters["browser"]; ok && browser != "" {
		prevWhereClause += " AND browser = ?"
		prevArgs = append(prevArgs, browser)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		prevWhereClause += " AND event_name = ?"
		prevArgs = append(prevArgs, eventName)
	}

	prevQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT session_id) as total_visits,
			COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_views
		FROM events 
		WHERE %s
	`, prevWhereClause)

	var prevTotalEvents, prevUniqueUsers, prevTotalVisits, prevPageViews int
	err = r.db.QueryRow(prevQuery, prevArgs...).Scan(&prevTotalEvents, &prevUniqueUsers, &prevTotalVisits, &prevPageViews)
	if err == nil {
		stats["prev_total_events"] = prevTotalEvents
		stats["prev_unique_users"] = prevUniqueUsers
		stats["prev_total_visits"] = prevTotalVisits
		stats["prev_page_views"] = prevPageViews

		// Calculate percentage changes
		if prevTotalEvents > 0 {
			stats["events_change"] = float64(totalEvents-prevTotalEvents) / float64(prevTotalEvents) * 100
		}
		if prevUniqueUsers > 0 {
			stats["users_change"] = float64(uniqueUsers-prevUniqueUsers) / float64(prevUniqueUsers) * 100
		}
		if prevTotalVisits > 0 {
			stats["visits_change"] = float64(totalVisits-prevTotalVisits) / float64(prevTotalVisits) * 100
		}
		if prevPageViews > 0 {
			stats["page_views_change"] = float64(pageViews-prevPageViews) / float64(prevPageViews) * 100
		}
	}

	return stats, nil
}

func (r *eventRepository) GetOnlineUsers(timeWindow int) (map[string]interface{}, error) {
	cutoffTime := time.Now().Add(-time.Duration(timeWindow) * time.Minute)

	query := `
		SELECT 
			COUNT(DISTINCT user_id) as online_users,
			COUNT(DISTINCT session_id) as active_sessions
		FROM events 
		WHERE timestamp >= ?
	`

	var onlineUsers, activeSessions int
	err := r.db.QueryRow(query, cutoffTime).Scan(&onlineUsers, &activeSessions)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"online_users":     onlineUsers,
		"active_sessions":  activeSessions,
		"time_window_mins": timeWindow,
		"cutoff_time":      cutoffTime,
	}, nil
}

func (r *eventRepository) GetProjects() ([]string, error) {
	query := `SELECT DISTINCT project_id FROM events WHERE project_id IS NOT NULL AND project_id != '' ORDER BY project_id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []string
	for rows.Next() {
		var projectID string
		if err := rows.Scan(&projectID); err != nil {
			continue
		}
		projects = append(projects, projectID)
	}

	return projects, nil
}
