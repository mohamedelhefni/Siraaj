package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

var migrations = []Migration{
	{
		Version:     1,
		Description: "Create events table",
		Up: `CREATE TABLE IF NOT EXISTS events (
			id UInt64,
			timestamp DateTime64(3),
			event_name LowCardinality(String),
			user_id String,
			session_id String,
			session_duration UInt32,
			url String,
			referrer String,
			user_agent String,
			ip String,
			country LowCardinality(String),
			browser LowCardinality(String),
			os LowCardinality(String),
			device LowCardinality(String),
			project_id LowCardinality(String) DEFAULT 'default',
			is_bot UInt8 DEFAULT 0,
			channel String DEFAULT ''
		) ENGINE = MergeTree()
		PARTITION BY toYYYYMM(timestamp)
		ORDER BY (project_id, timestamp, event_name, is_bot)
		SETTINGS index_granularity = 8192`,
		Down: `DROP TABLE IF EXISTS events`,
	},
	{
		Version:     2,
		Description: "Add MATERIALIZED date columns for faster date-based queries",
		Up: `ALTER TABLE events 
			ADD COLUMN IF NOT EXISTS date_hour DateTime MATERIALIZED toStartOfHour(timestamp),
			ADD COLUMN IF NOT EXISTS date_day Date MATERIALIZED toDate(timestamp),
			ADD COLUMN IF NOT EXISTS date_month Date MATERIALIZED toStartOfMonth(timestamp)`,
		Down: `ALTER TABLE events 
			DROP COLUMN IF EXISTS date_hour,
			DROP COLUMN IF EXISTS date_day,
			DROP COLUMN IF EXISTS date_month`,
	},
	{
		Version:     3,
		Description: "Add projections for common aggregation queries",
		Up: `ALTER TABLE events ADD PROJECTION IF NOT EXISTS stats_by_country (
			SELECT 
				country,
				toStartOfHour(timestamp) as hour,
				COUNT() as event_count,
				uniq(user_id) as user_count,
				uniq(session_id) as session_count
			GROUP BY country, hour
		)`,
		Down: `ALTER TABLE events DROP PROJECTION IF EXISTS stats_by_country`,
	},
	{
		Version:     4,
		Description: "Add projection for browser/device/OS stats",
		Up: `ALTER TABLE events ADD PROJECTION IF NOT EXISTS stats_by_device (
			SELECT 
				browser,
				device,
				os,
				toDate(timestamp) as day,
				COUNT() as event_count,
				uniq(user_id) as user_count
			GROUP BY browser, device, os, day
		)`,
		Down: `ALTER TABLE events DROP PROJECTION IF EXISTS stats_by_device`,
	},
	{
		Version:     5,
		Description: "Add projection for page stats",
		Up: `ALTER TABLE events ADD PROJECTION IF NOT EXISTS stats_by_page (
			SELECT 
				url,
				toDate(timestamp) as day,
				COUNT() as event_count,
				uniq(session_id) as session_count
			GROUP BY url, day
		)`,
		Down: `ALTER TABLE events DROP PROJECTION IF EXISTS stats_by_page`,
	}, {
		Version:     6,
		Description: "Create daily aggregated stats materialized view",
		Up: `
CREATE MATERIALIZED VIEW events_daily_stats
ENGINE = AggregatingMergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (date, project_id)
POPULATE AS
SELECT
    toDate(timestamp) as date,
    project_id,
    
    -- Basic counts
    countState() as total_events_state,
    uniqState(user_id) as unique_users_state,
    uniqState(session_id) as total_visits_state,
    
    -- Page view metrics
    countStateIf(event_name = 'page_view') as page_views_state,
    uniqStateIf(session_id, event_name = 'page_view') as sessions_with_views_state,
    
    -- Session duration
    avgStateIf(session_duration, session_duration > 0) as avg_session_duration_state,
    
    -- Bot metrics
    countStateIf(is_bot = 1) as bot_events_state,
    countStateIf(is_bot = 0) as human_events_state,
    uniqStateIf(user_id, is_bot = 1) as bot_users_state,
    uniqStateIf(user_id, is_bot = 0) as human_users_state
    
FROM events
GROUP BY date, project_id;
		`,
		Down: `DROP MATERIALIZED VIEW IF EXISTS events_daily_stats`,
	},
}

func initMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version UInt32,
			description String,
			applied_at DateTime
		) ENGINE = MergeTree()
		ORDER BY version
	`)
	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT max(version) FROM schema_migrations").Scan(&version)
	if err != nil {
		// If table is empty, return 0
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

func Migrate(db *sql.DB) error {
	log.Println("Running database migrations...")

	if err := initMigrationTable(db); err != nil {
		return fmt.Errorf("failed to initialize migration table: %v", err)
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	log.Printf("Current schema version: %d", currentVersion)

	for _, m := range migrations {
		if m.Version <= currentVersion {
			continue
		}

		log.Printf("Applying migration %d: %s", m.Version, m.Description)

		if _, err := db.Exec(m.Up); err != nil {
			return fmt.Errorf("failed to apply migration %d: %v", m.Version, err)
		}

		if _, err := db.Exec(
			"INSERT INTO schema_migrations (version, description, applied_at) VALUES (?, ?, ?)",
			m.Version, m.Description, time.Now(),
		); err != nil {
			return fmt.Errorf("failed to record migration %d: %v", m.Version, err)
		}

		log.Printf("✓ Successfully applied migration %d", m.Version)
	}

	log.Println("✓ All migrations completed")
	return nil
}

func Rollback(db *sql.DB, targetVersion int) error {
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	if targetVersion >= currentVersion {
		return fmt.Errorf("target version must be less than current version")
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version <= targetVersion || m.Version > currentVersion {
			continue
		}

		log.Printf("Rolling back migration %d: %s", m.Version, m.Description)

		if _, err := db.Exec(m.Down); err != nil {
			return fmt.Errorf("failed to rollback migration %d: %v", m.Version, err)
		}

		if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = ?", m.Version); err != nil {
			return fmt.Errorf("failed to remove migration record %d: %v", m.Version, err)
		}

		log.Printf("✓ Successfully rolled back migration %d", m.Version)
	}

	return nil
}
