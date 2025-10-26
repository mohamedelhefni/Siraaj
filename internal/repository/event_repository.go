package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	GetFunnelAnalysis(request domain.FunnelRequest) (*domain.FunnelAnalysisResult, error)
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event domain.Event) error {
	query := `
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration, url, referrer, 
			user_agent, ip, country, browser, os, device, is_bot, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if event.ProjectID == "" {
		event.ProjectID = "default"
	}

	_, err := r.db.Exec(query,
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
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO events (id, timestamp, event_name, user_id, session_id, session_duration, url, referrer, 
			user_agent, ip, country, browser, os, device, is_bot, project_id)
		VALUES (nextval('id_sequence'), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: failed to close statement: %v", err)
		}
	}()

	for _, event := range events {
		if event.ProjectID == "" {
			event.ProjectID = "default"
		}

		_, err := stmt.Exec(
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
		SELECT id, timestamp, event_name, user_id, session_id, session_duration, url, referrer,
			user_agent, ip, country, browser, os, device, is_bot, project_id
		FROM events
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Warning: failed to close rows: %v", err)
		}
	}()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		err := rows.Scan(
			&e.ID, &e.Timestamp, &e.EventName, &e.UserID, &e.SessionID, &e.SessionDuration,
			&e.URL, &e.Referrer, &e.UserAgent, &e.IP, &e.Country,
			&e.Browser, &e.OS, &e.Device, &e.IsBot, &e.ProjectID,
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
	if device, ok := filters["device"]; ok && device != "" {
		whereClause += " AND device = ?"
		args = append(args, device)
	}
	if os, ok := filters["os"]; ok && os != "" {
		whereClause += " AND os = ?"
		args = append(args, os)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		whereClause += " AND event_name = ?"
		args = append(args, eventName)
	}
	if page, ok := filters["page"]; ok && page != "" {
		whereClause += " AND url = ?"
		args = append(args, page)
	}
	if botFilter, ok := filters["botFilter"]; ok && botFilter != "" {
		if botFilter == "bot" {
			whereClause += " AND is_bot = TRUE"
		} else if botFilter == "human" {
			whereClause += " AND is_bot = FALSE"
		}
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
			COUNT(DISTINCT CASE WHEN event_name = 'page_view' THEN session_id END) as sessions_with_views,
			AVG(CASE WHEN session_duration > 0 THEN session_duration END) as avg_session_duration
		FROM date_filtered
	)
	SELECT total_events, unique_users, total_visits, page_views, sessions_with_views, avg_session_duration FROM event_stats;
	`, whereClause)

	var totalEvents, uniqueUsers, totalVisits, pageViews, sessionsWithViews int
	var avgSessionDuration sql.NullFloat64
	err := r.db.QueryRow(aggregateQuery, args...).Scan(&totalEvents, &uniqueUsers, &totalVisits, &pageViews, &sessionsWithViews, &avgSessionDuration)
	if err != nil {
		return nil, err
	}
	stats["total_events"] = totalEvents
	stats["unique_users"] = uniqueUsers
	stats["total_visits"] = totalVisits
	stats["page_views"] = pageViews

	// Add average session duration (default to 0 if NULL)
	if avgSessionDuration.Valid {
		stats["avg_session_duration"] = avgSessionDuration.Float64
	} else {
		stats["avg_session_duration"] = 0.0
	}

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

	// Bot statistics
	botQuery := fmt.Sprintf(`
		SELECT 
			COUNT(CASE WHEN is_bot = TRUE THEN 1 END) as bot_events,
			COUNT(CASE WHEN is_bot = FALSE THEN 1 END) as human_events,
			COUNT(DISTINCT CASE WHEN is_bot = TRUE THEN user_id END) as bot_users,
			COUNT(DISTINCT CASE WHEN is_bot = FALSE THEN user_id END) as human_users
		FROM events 
		WHERE %s
	`, whereClause)

	var botEvents, humanEvents, botUsers, humanUsers int
	err = r.db.QueryRow(botQuery, args...).Scan(&botEvents, &humanEvents, &botUsers, &humanUsers)
	if err == nil {
		stats["bot_events"] = botEvents
		stats["human_events"] = humanEvents
		stats["bot_users"] = botUsers
		stats["human_users"] = humanUsers

		// Calculate bot percentage
		if totalEvents > 0 {
			stats["bot_percentage"] = float64(botEvents) / float64(totalEvents) * 100
		} else {
			stats["bot_percentage"] = 0.0
		}
	}

	// Top Events with optimized query
	query := fmt.Sprintf(`
		SELECT event_name, COUNT(*) as count 
		FROM events 
		WHERE %s
		GROUP BY event_name 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)
	queryArgs := append(args, limit)

	topEventsRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if topEventsRows != nil {
			if err := topEventsRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	topEvents := []map[string]interface{}{}
	for topEventsRows.Next() {
		var name string
		var count int
		if err := topEventsRows.Scan(&name, &count); err != nil {
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

	// Determine what metric to display in timeline
	metric := filters["metric"]
	var selectClause string
	switch metric {
	case "users":
		selectClause = "COUNT(DISTINCT user_id) as count"
	case "visits":
		selectClause = "COUNT(DISTINCT session_id) as count"
	case "page_views":
		selectClause = "COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as count"
	case "events":
		selectClause = "COUNT(*) as count"
	case "views_per_visit":
		selectClause = "CAST(COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) AS FLOAT) / NULLIF(COUNT(DISTINCT session_id), 0) as count"
	case "bounce_rate":
		// For bounce rate in timeline, we need to use a different approach
		// We'll calculate it per time period using a window function or aggregation
		// This is a simplified version that's much faster
		selectClause = `
			CASE 
				WHEN COUNT(DISTINCT session_id) = 0 THEN 0
				ELSE CAST(SUM(CASE WHEN event_name = 'page_view' THEN 1 ELSE 0 END) AS FLOAT) * 100.0 / NULLIF(COUNT(DISTINCT session_id), 0)
			END as count`
	case "visit_duration":
		selectClause = "AVG(CASE WHEN session_duration > 0 THEN session_duration END) as count"
	default: // Default to users
		selectClause = "COUNT(DISTINCT user_id) as count"
	}

	// Determine granularity based on date range
	if timelineDuration <= 24*time.Hour {
		// For today or single day: show hourly data
		if metric == "bounce_rate" {
			// Special optimized query for bounce rate
			timelineQuery = fmt.Sprintf(`
				WITH session_page_counts AS (
					SELECT 
						strftime(DATE_TRUNC('hour', timestamp), '%%Y-%%m-%%d %%H:00:00') as date,
						session_id,
						COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_view_count
					FROM events 
					WHERE %s
					GROUP BY date, session_id
				)
				SELECT 
					date,
					CAST(COUNT(CASE WHEN page_view_count = 1 THEN 1 END) AS FLOAT) * 100.0 / NULLIF(COUNT(*), 0) as count
				FROM session_page_counts
				GROUP BY date
				ORDER BY date
			`, whereClause)
		} else {
			timelineQuery = fmt.Sprintf(`
				SELECT 
					strftime(DATE_TRUNC('hour', timestamp), '%%Y-%%m-%%d %%H:00:00') as date, 
					%s
				FROM events 
				WHERE %s
				GROUP BY date 
				ORDER BY date
			`, selectClause, whereClause)
		}
		timeFormat = "hour"
	} else if timelineDuration <= 90*24*time.Hour {
		// For up to 3 months: show daily data
		if metric == "bounce_rate" {
			// Special optimized query for bounce rate
			timelineQuery = fmt.Sprintf(`
				WITH session_page_counts AS (
					SELECT 
						strftime(DATE_TRUNC('day', timestamp), '%%Y-%%m-%%d') as date,
						session_id,
						COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_view_count
					FROM events 
					WHERE %s
					GROUP BY date, session_id
				)
				SELECT 
					date,
					CAST(COUNT(CASE WHEN page_view_count = 1 THEN 1 END) AS FLOAT) * 100.0 / NULLIF(COUNT(*), 0) as count
				FROM session_page_counts
				GROUP BY date
				ORDER BY date
			`, whereClause)
		} else {
			timelineQuery = fmt.Sprintf(`
				SELECT 
					strftime(DATE_TRUNC('day', timestamp), '%%Y-%%m-%%d') as date, 
					%s
				FROM events 
				WHERE %s
				GROUP BY date 
				ORDER BY date
			`, selectClause, whereClause)
		}
		timeFormat = "day"
	} else {
		// For more than 3 months: show monthly data
		if metric == "bounce_rate" {
			// Special optimized query for bounce rate
			timelineQuery = fmt.Sprintf(`
				WITH session_page_counts AS (
					SELECT 
						strftime(DATE_TRUNC('month', timestamp), '%%Y-%%m-01') as date,
						session_id,
						COUNT(CASE WHEN event_name = 'page_view' THEN 1 END) as page_view_count
					FROM events 
					WHERE %s
					GROUP BY date, session_id
				)
				SELECT 
					date,
					CAST(COUNT(CASE WHEN page_view_count = 1 THEN 1 END) AS FLOAT) * 100.0 / NULLIF(COUNT(*), 0) as count
				FROM session_page_counts
				GROUP BY date
				ORDER BY date
			`, whereClause)
		} else {
			timelineQuery = fmt.Sprintf(`
				SELECT 
					strftime(DATE_TRUNC('month', timestamp), '%%Y-%%m-01') as date, 
					%s
				FROM events 
				WHERE %s
				GROUP BY date 
				ORDER BY date
			`, selectClause, whereClause)
		}
		timeFormat = "month"
	}

	timelineRows, err := r.db.Query(timelineQuery, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if timelineRows != nil {
			if err := timelineRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	timeline := []map[string]interface{}{}
	for timelineRows.Next() {
		var date string
		var count sql.NullFloat64
		if err := timelineRows.Scan(&date, &count); err != nil {
			log.Printf("Error scanning timeline row: %v", err)
			continue
		}

		// Use float64 value if valid, otherwise 0
		countValue := 0.0
		if count.Valid {
			countValue = count.Float64
		}

		timeline = append(timeline, map[string]interface{}{
			"date":  date,
			"count": countValue,
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

	topPagesRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if topPagesRows != nil {
			if err := topPagesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	topPages := []map[string]interface{}{}
	for topPagesRows.Next() {
		var url string
		var count int
		if err := topPagesRows.Scan(&url, &count); err != nil {
			continue
		}
		topPages = append(topPages, map[string]interface{}{
			"url":   url,
			"count": count,
		})
	}
	stats["top_pages"] = topPages

	// Entry Pages (first page in each session)
	entryPagesQuery := fmt.Sprintf(`
		WITH entry_pages AS (
			SELECT DISTINCT ON (session_id) 
				session_id, 
				url
			FROM events 
			WHERE %s AND event_name = 'page_view' AND url IS NOT NULL AND url != ''
			ORDER BY session_id, timestamp ASC
		)
		SELECT url, COUNT(*) as count
		FROM entry_pages
		GROUP BY url
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	entryPagesRows, err := r.db.Query(entryPagesQuery, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if entryPagesRows != nil {
			if err := entryPagesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	entryPages := []map[string]interface{}{}
	for entryPagesRows.Next() {
		var url string
		var count int
		if err := entryPagesRows.Scan(&url, &count); err != nil {
			continue
		}
		entryPages = append(entryPages, map[string]interface{}{
			"url":   url,
			"count": count,
		})
	}
	stats["entry_pages"] = entryPages

	// Exit Pages (last page in each session)
	exitPagesQuery := fmt.Sprintf(`
		WITH exit_pages AS (
			SELECT DISTINCT ON (session_id) 
				session_id, 
				url
			FROM events 
			WHERE %s AND event_name = 'page_view' AND url IS NOT NULL AND url != ''
			ORDER BY session_id, timestamp DESC
		)
		SELECT url, COUNT(*) as count
		FROM exit_pages
		GROUP BY url
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	exitPagesRows, err := r.db.Query(exitPagesQuery, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if exitPagesRows != nil {
			if err := exitPagesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	exitPages := []map[string]interface{}{}
	for exitPagesRows.Next() {
		var url string
		var count int
		if err := exitPagesRows.Scan(&url, &count); err != nil {
			continue
		}
		exitPages = append(exitPages, map[string]interface{}{
			"url":   url,
			"count": count,
		})
	}
	stats["exit_pages"] = exitPages

	// Browsers
	query = fmt.Sprintf(`
		SELECT browser, COUNT(*) as count 
		FROM events 
		WHERE %s AND browser IS NOT NULL AND browser != ''
		GROUP BY browser 
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	browsersRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if browsersRows != nil {
			if err := browsersRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	browsers := []map[string]interface{}{}
	for browsersRows.Next() {
		var browser string
		var count int
		if err := browsersRows.Scan(&browser, &count); err != nil {
			continue
		}
		browsers = append(browsers, map[string]interface{}{
			"name":  browser,
			"count": count,
		})
	}
	stats["browsers"] = browsers

	// Devices
	query = fmt.Sprintf(`
		SELECT device, COUNT(*) as count 
		FROM events 
		WHERE %s AND device IS NOT NULL AND device != ''
		GROUP BY device 
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	devicesRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if devicesRows != nil {
			if err := devicesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	devices := []map[string]interface{}{}
	for devicesRows.Next() {
		var device string
		var count int
		if err := devicesRows.Scan(&device, &count); err != nil {
			continue
		}
		devices = append(devices, map[string]interface{}{
			"name":  device,
			"count": count,
		})
	}
	stats["devices"] = devices

	// Operating Systems
	query = fmt.Sprintf(`
		SELECT os, COUNT(*) as count 
		FROM events 
		WHERE %s AND os IS NOT NULL AND os != ''
		GROUP BY os 
		ORDER BY count DESC
		LIMIT ?
	`, whereClause)

	osRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if osRows != nil {
			if err := osRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	operatingSystems := []map[string]interface{}{}
	for osRows.Next() {
		var os string
		var count int
		if err := osRows.Scan(&os, &count); err != nil {
			continue
		}
		operatingSystems = append(operatingSystems, map[string]interface{}{
			"name":  os,
			"count": count,
		})
	}
	stats["os"] = operatingSystems

	// Top Countries
	query = fmt.Sprintf(`
		SELECT country, COUNT(*) as count 
		FROM events 
		WHERE %s AND country IS NOT NULL AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT ?
	`, whereClause)

	countriesRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if countriesRows != nil {
			if err := countriesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	topCountries := []map[string]interface{}{}
	for countriesRows.Next() {
		var country string
		var count int
		if err := countriesRows.Scan(&country, &count); err != nil {
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

	sourcesRows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if sourcesRows != nil {
			if err := sourcesRows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}
	}()

	topSources := []map[string]interface{}{}
	for sourcesRows.Next() {
		var referrer string
		var count int
		if err := sourcesRows.Scan(&referrer, &count); err != nil {
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
	if device, ok := filters["device"]; ok && device != "" {
		prevWhereClause += " AND device = ?"
		prevArgs = append(prevArgs, device)
	}
	if os, ok := filters["os"]; ok && os != "" {
		prevWhereClause += " AND os = ?"
		prevArgs = append(prevArgs, os)
	}
	if eventName, ok := filters["event"]; ok && eventName != "" {
		prevWhereClause += " AND event_name = ?"
		prevArgs = append(prevArgs, eventName)
	}
	if page, ok := filters["page"]; ok && page != "" {
		prevWhereClause += " AND url = ?"
		prevArgs = append(prevArgs, page)
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
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Warning: failed to close rows: %v", err)
		}
	}()

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

func (r *eventRepository) GetFunnelAnalysis(request domain.FunnelRequest) (*domain.FunnelAnalysisResult, error) {
	if len(request.Steps) == 0 {
		return nil, fmt.Errorf("at least one funnel step is required")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// Set to beginning and end of day
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	result := &domain.FunnelAnalysisResult{
		Steps:     make([]domain.FunnelStepResult, len(request.Steps)),
		TimeRange: fmt.Sprintf("%s to %s", request.StartDate, request.EndDate),
	}

	// Build base WHERE clause for global filters
	baseWhereClause := "timestamp BETWEEN ? AND ?"
	baseArgs := []interface{}{startDate, endDate}

	if projectID, ok := request.Filters["project"]; ok && projectID != "" {
		baseWhereClause += " AND project_id = ?"
		baseArgs = append(baseArgs, projectID)
	}
	if country, ok := request.Filters["country"]; ok && country != "" {
		baseWhereClause += " AND country = ?"
		baseArgs = append(baseArgs, country)
	}
	if browser, ok := request.Filters["browser"]; ok && browser != "" {
		baseWhereClause += " AND browser = ?"
		baseArgs = append(baseArgs, browser)
	}
	if device, ok := request.Filters["device"]; ok && device != "" {
		baseWhereClause += " AND device = ?"
		baseArgs = append(baseArgs, device)
	}
	if os, ok := request.Filters["os"]; ok && os != "" {
		baseWhereClause += " AND os = ?"
		baseArgs = append(baseArgs, os)
	}
	if botFilter, ok := request.Filters["botFilter"]; ok && botFilter != "" {
		if botFilter == "bot" {
			baseWhereClause += " AND is_bot = TRUE"
		} else if botFilter == "human" {
			baseWhereClause += " AND is_bot = FALSE"
		}
	}

	// For each step, calculate metrics
	var previousUserCount int64 = 0
	var totalUsers int64 = 0

	for i, step := range request.Steps {
		// Build WHERE clause for this step
		stepWhereClause := baseWhereClause
		stepArgs := make([]interface{}, len(baseArgs))
		copy(stepArgs, baseArgs)

		// Add event name filter
		if step.EventName != "" {
			stepWhereClause += " AND event_name = ?"
			stepArgs = append(stepArgs, step.EventName)
		}

		// Add URL filter if specified
		if step.URL != "" {
			stepWhereClause += " AND url = ?"
			stepArgs = append(stepArgs, step.URL)
		}

		// Add step-specific filters
		for key, value := range step.Filters {
			switch key {
			case "country":
				stepWhereClause += " AND country = ?"
				stepArgs = append(stepArgs, value)
			case "browser":
				stepWhereClause += " AND browser = ?"
				stepArgs = append(stepArgs, value)
			case "device":
				stepWhereClause += " AND device = ?"
				stepArgs = append(stepArgs, value)
			case "os":
				stepWhereClause += " AND os = ?"
				stepArgs = append(stepArgs, value)
			}
		}

		// If this is not the first step, we need to filter for users who completed previous steps
		if i == 0 {
			// First step: count all matching users
			query := fmt.Sprintf(`
				SELECT 
					COUNT(DISTINCT user_id) as user_count,
					COUNT(DISTINCT session_id) as session_count,
					COUNT(*) as event_count
				FROM events 
				WHERE %s
			`, stepWhereClause)

			var userCount, sessionCount, eventCount int64
			err := r.db.QueryRow(query, stepArgs...).Scan(&userCount, &sessionCount, &eventCount)
			if err != nil {
				return nil, fmt.Errorf("error querying step %d: %v", i+1, err)
			}

			result.Steps[i] = domain.FunnelStepResult{
				Step:           step,
				UserCount:      userCount,
				SessionCount:   sessionCount,
				EventCount:     eventCount,
				ConversionRate: 100.0, // First step is always 100%
				OverallRate:    100.0,
				DropoffRate:    0.0,
			}

			totalUsers = userCount
			previousUserCount = userCount
			result.TotalUsers = totalUsers

		} else {
			// Subsequent steps: only count users who completed all previous steps
			// Build a CTE that finds users who completed all previous steps in order
			var cteBuilder strings.Builder
			cteBuilder.WriteString("WITH ")

			// Collect all arguments for all CTEs
			var allCteArgs []interface{}

			// Create CTEs for each previous step
			for j := 0; j <= i; j++ {
				if j > 0 {
					cteBuilder.WriteString(", ")
				}

				prevStep := request.Steps[j]
				cteName := fmt.Sprintf("step_%d", j+1)

				// Build WHERE for this CTE
				var cteWhereClause string
				var cteArgs []interface{}

				if j == 0 {
					// First step: simple query without joins
					cteWhereClause = baseWhereClause
					cteArgs = make([]interface{}, len(baseArgs))
					copy(cteArgs, baseArgs)

					if prevStep.EventName != "" {
						cteWhereClause += " AND event_name = ?"
						cteArgs = append(cteArgs, prevStep.EventName)
					}
					if prevStep.URL != "" {
						cteWhereClause += " AND url = ?"
						cteArgs = append(cteArgs, prevStep.URL)
					}

					for key, value := range prevStep.Filters {
						switch key {
						case "country":
							cteWhereClause += " AND country = ?"
							cteArgs = append(cteArgs, value)
						case "browser":
							cteWhereClause += " AND browser = ?"
							cteArgs = append(cteArgs, value)
						case "device":
							cteWhereClause += " AND device = ?"
							cteArgs = append(cteArgs, value)
						case "os":
							cteWhereClause += " AND os = ?"
							cteArgs = append(cteArgs, value)
						}
					}

					fmt.Fprintf(&cteBuilder, "%s AS (SELECT user_id, session_id, timestamp FROM events WHERE %s)", cteName, cteWhereClause)
					allCteArgs = append(allCteArgs, cteArgs...)
				} else {
					// Subsequent steps: join with previous step
					// Build WHERE clause with e. prefix
					cteWhereClause = "e.timestamp BETWEEN ? AND ?"
					cteArgs = []interface{}{startDate, endDate}

					// Add global filters with e. prefix
					if projectID, ok := request.Filters["project"]; ok && projectID != "" {
						cteWhereClause += " AND e.project_id = ?"
						cteArgs = append(cteArgs, projectID)
					}
					if country, ok := request.Filters["country"]; ok && country != "" {
						cteWhereClause += " AND e.country = ?"
						cteArgs = append(cteArgs, country)
					}
					if browser, ok := request.Filters["browser"]; ok && browser != "" {
						cteWhereClause += " AND e.browser = ?"
						cteArgs = append(cteArgs, browser)
					}
					if device, ok := request.Filters["device"]; ok && device != "" {
						cteWhereClause += " AND e.device = ?"
						cteArgs = append(cteArgs, device)
					}
					if os, ok := request.Filters["os"]; ok && os != "" {
						cteWhereClause += " AND e.os = ?"
						cteArgs = append(cteArgs, os)
					}
					if botFilter, ok := request.Filters["botFilter"]; ok && botFilter != "" {
						if botFilter == "bot" {
							cteWhereClause += " AND e.is_bot = TRUE"
						} else if botFilter == "human" {
							cteWhereClause += " AND e.is_bot = FALSE"
						}
					}

					if prevStep.EventName != "" {
						cteWhereClause += " AND e.event_name = ?"
						cteArgs = append(cteArgs, prevStep.EventName)
					}
					if prevStep.URL != "" {
						cteWhereClause += " AND e.url = ?"
						cteArgs = append(cteArgs, prevStep.URL)
					}

					for key, value := range prevStep.Filters {
						switch key {
						case "country":
							cteWhereClause += " AND e.country = ?"
							cteArgs = append(cteArgs, value)
						case "browser":
							cteWhereClause += " AND e.browser = ?"
							cteArgs = append(cteArgs, value)
						case "device":
							cteWhereClause += " AND e.device = ?"
							cteArgs = append(cteArgs, value)
						case "os":
							cteWhereClause += " AND e.os = ?"
							cteArgs = append(cteArgs, value)
						}
					}

					prevCteName := fmt.Sprintf("step_%d", j)
					fmt.Fprintf(&cteBuilder, "%s AS (SELECT e.user_id, e.session_id, e.timestamp FROM events e INNER JOIN %s prev ON e.user_id = prev.user_id AND e.timestamp > prev.timestamp WHERE %s)", cteName, prevCteName, cteWhereClause)
					allCteArgs = append(allCteArgs, cteArgs...)
				}
			}

			// Main query to count users who reached this step
			currentCteName := fmt.Sprintf("step_%d", i+1)
			mainQuery := fmt.Sprintf(`
				%s
				SELECT 
					COUNT(DISTINCT user_id) as user_count,
					COUNT(DISTINCT session_id) as session_count,
					COUNT(*) as event_count
				FROM %s
			`, cteBuilder.String(), currentCteName)

			var userCount, sessionCount, eventCount int64
			err := r.db.QueryRow(mainQuery, allCteArgs...).Scan(&userCount, &sessionCount, &eventCount)
			if err != nil {
				return nil, fmt.Errorf("error querying step %d: %v", i+1, err)
			}

			// Calculate conversion rates
			conversionRate := 0.0
			if previousUserCount > 0 {
				conversionRate = float64(userCount) / float64(previousUserCount) * 100
			}

			overallRate := 0.0
			if totalUsers > 0 {
				overallRate = float64(userCount) / float64(totalUsers) * 100
			}

			dropoffRate := 100.0 - conversionRate

			result.Steps[i] = domain.FunnelStepResult{
				Step:           step,
				UserCount:      userCount,
				SessionCount:   sessionCount,
				EventCount:     eventCount,
				ConversionRate: conversionRate,
				OverallRate:    overallRate,
				DropoffRate:    dropoffRate,
			}

			previousUserCount = userCount
		}

		// Calculate average and median time to next step (if not the last step)
		if i < len(request.Steps)-1 {
			nextStep := request.Steps[i+1]

			// Build query to find time between this step and next step
			timeQuery := fmt.Sprintf(`
				WITH current_step AS (
					SELECT user_id, timestamp 
					FROM events 
					WHERE %s
				),
				next_step AS (
					SELECT user_id, timestamp 
					FROM events 
					WHERE %s
				)
				SELECT 
					AVG(EXTRACT(EPOCH FROM (n.timestamp - c.timestamp))) as avg_time,
					PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (n.timestamp - c.timestamp))) as median_time
				FROM current_step c
				INNER JOIN next_step n ON c.user_id = n.user_id AND n.timestamp > c.timestamp
			`, stepWhereClause, stepWhereClause) // We'll need to build next step WHERE clause

			// Build next step WHERE clause
			nextStepWhereClause := baseWhereClause
			nextStepArgs := make([]interface{}, len(baseArgs))
			copy(nextStepArgs, baseArgs)

			if nextStep.EventName != "" {
				nextStepWhereClause += " AND event_name = ?"
				nextStepArgs = append(nextStepArgs, nextStep.EventName)
			}
			if nextStep.URL != "" {
				nextStepWhereClause += " AND url = ?"
				nextStepArgs = append(nextStepArgs, nextStep.URL)
			}

			// Combine args for the time query
			timeQueryArgs := append(stepArgs, nextStepArgs...)

			var avgTime, medianTime sql.NullFloat64
			err := r.db.QueryRow(timeQuery, timeQueryArgs...).Scan(&avgTime, &medianTime)
			if err == nil {
				if avgTime.Valid {
					result.Steps[i].AvgTimeToNext = avgTime.Float64
				}
				if medianTime.Valid {
					result.Steps[i].MedianTimeToNext = medianTime.Float64
				}
			}
		}
	}

	// Calculate overall completion metrics
	if len(request.Steps) > 0 {
		lastStep := result.Steps[len(result.Steps)-1]
		result.CompletedUsers = lastStep.UserCount

		if result.TotalUsers > 0 {
			result.CompletionRate = float64(result.CompletedUsers) / float64(result.TotalUsers) * 100
		}

		// Calculate average time to complete entire funnel
		if len(request.Steps) > 1 {
			firstStep := request.Steps[0]
			lastStepDef := request.Steps[len(request.Steps)-1]

			// Build WHERE clauses
			firstWhereClause := baseWhereClause
			firstArgs := make([]interface{}, len(baseArgs))
			copy(firstArgs, baseArgs)
			if firstStep.EventName != "" {
				firstWhereClause += " AND event_name = ?"
				firstArgs = append(firstArgs, firstStep.EventName)
			}
			if firstStep.URL != "" {
				firstWhereClause += " AND url = ?"
				firstArgs = append(firstArgs, firstStep.URL)
			}

			lastWhereClause := baseWhereClause
			lastArgs := make([]interface{}, len(baseArgs))
			copy(lastArgs, baseArgs)
			if lastStepDef.EventName != "" {
				lastWhereClause += " AND event_name = ?"
				lastArgs = append(lastArgs, lastStepDef.EventName)
			}
			if lastStepDef.URL != "" {
				lastWhereClause += " AND url = ?"
				lastArgs = append(lastArgs, lastStepDef.URL)
			}

			completionTimeQuery := fmt.Sprintf(`
				WITH first_step AS (
					SELECT user_id, MIN(timestamp) as first_time 
					FROM events 
					WHERE %s
					GROUP BY user_id
				),
				last_step AS (
					SELECT user_id, MAX(timestamp) as last_time 
					FROM events 
					WHERE %s
					GROUP BY user_id
				)
				SELECT AVG(EXTRACT(EPOCH FROM (l.last_time - f.first_time))) as avg_completion
				FROM first_step f
				INNER JOIN last_step l ON f.user_id = l.user_id AND l.last_time > f.first_time
			`, firstWhereClause, lastWhereClause)

			completionArgs := append(firstArgs, lastArgs...)

			var avgCompletion sql.NullFloat64
			err := r.db.QueryRow(completionTimeQuery, completionArgs...).Scan(&avgCompletion)
			if err == nil && avgCompletion.Valid {
				result.AvgCompletion = avgCompletion.Float64
			}
		}
	}

	return result, nil
}
