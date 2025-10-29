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
			id UBIGINT PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			event_name VARCHAR NOT NULL,
			user_id VARCHAR,
			session_id VARCHAR,
			session_duration INTEGER,
			url VARCHAR,
			referrer VARCHAR,
			user_agent VARCHAR,
			ip VARCHAR,
			country VARCHAR,
			browser VARCHAR,
			os VARCHAR,
			device VARCHAR,
			project_id VARCHAR DEFAULT 'default'
		)`,
		Down: `DROP TABLE IF EXISTS events`,
	},
	{
		Version:     2,
		Description: "Create id sequence",
		Up:          `CREATE SEQUENCE IF NOT EXISTS id_sequence START 1`,
		Down:        `DROP SEQUENCE IF EXISTS id_sequence`,
	},
	{
		Version:     3,
		Description: "Create indexes",
		Up: `CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_event_name ON events(event_name);
		CREATE INDEX IF NOT EXISTS idx_user_id ON events(user_id);
		CREATE INDEX IF NOT EXISTS idx_country ON events(country);
		CREATE INDEX IF NOT EXISTS idx_referrer ON events(referrer);
		CREATE INDEX IF NOT EXISTS idx_project_id ON events(project_id);
		CREATE INDEX IF NOT EXISTS idx_session_id ON events(session_id)`,
		Down: `DROP INDEX IF EXISTS idx_timestamp;
		DROP INDEX IF EXISTS idx_event_name;
		DROP INDEX IF EXISTS idx_user_id;
		DROP INDEX IF EXISTS idx_country;
		DROP INDEX IF EXISTS idx_referrer;
		DROP INDEX IF EXISTS idx_project_id;
		DROP INDEX IF EXISTS idx_session_id`,
	},
	{
		Version:     4,
		Description: "Add is_bot column to events table",
		Up:          `ALTER TABLE events ADD COLUMN IF NOT EXISTS is_bot BOOLEAN DEFAULT FALSE`,
		Down:        `ALTER TABLE events DROP COLUMN IF EXISTS is_bot`,
	},
	{
		Version:     5,
		Description: "Create index on is_bot column",
		Up:          `CREATE INDEX IF NOT EXISTS idx_is_bot ON events(is_bot)`,
		Down:        `DROP INDEX IF EXISTS idx_is_bot`,
	},
	{
		Version:     6,
		Description: "Create composite indexes for common query patterns",
		Up: `-- Composite index for timeline queries (most common query pattern)
		CREATE INDEX IF NOT EXISTS idx_timestamp_event_name ON events(timestamp DESC, event_name);
		
		-- Composite index for filtered queries
		CREATE INDEX IF NOT EXISTS idx_timestamp_project ON events(timestamp DESC, project_id);
		
		-- Composite index for session analysis
		CREATE INDEX IF NOT EXISTS idx_session_timestamp ON events(session_id, timestamp);
		
		-- Composite index for user journey analysis
		CREATE INDEX IF NOT EXISTS idx_user_timestamp ON events(user_id, timestamp);
		
		-- Composite index for URL-based queries
		CREATE INDEX IF NOT EXISTS idx_timestamp_url ON events(timestamp DESC, url);
		
		-- Composite index for country filtering
		CREATE INDEX IF NOT EXISTS idx_timestamp_country ON events(timestamp DESC, country);`,
		Down: `DROP INDEX IF EXISTS idx_timestamp_event_name;
		DROP INDEX IF EXISTS idx_timestamp_project;
		DROP INDEX IF EXISTS idx_session_timestamp;
		DROP INDEX IF EXISTS idx_user_timestamp;
		DROP INDEX IF EXISTS idx_timestamp_url;
		DROP INDEX IF EXISTS idx_timestamp_country`,
	},
	{
		Version:     7,
		Description: "Add additional performance indexes for analytics queries",
		Up: `-- Composite index for bot filtering with timestamp
		CREATE INDEX IF NOT EXISTS idx_timestamp_is_bot ON events(timestamp DESC, is_bot);
		
		-- Composite index for browser analytics
		CREATE INDEX IF NOT EXISTS idx_timestamp_browser ON events(timestamp DESC, browser);
		
		-- Composite index for device analytics
		CREATE INDEX IF NOT EXISTS idx_timestamp_device ON events(timestamp DESC, device);
		
		-- Composite index for OS analytics
		CREATE INDEX IF NOT EXISTS idx_timestamp_os ON events(timestamp DESC, os);
		
		-- Composite index for referrer/source analytics
		CREATE INDEX IF NOT EXISTS idx_timestamp_referrer ON events(timestamp DESC, referrer);
		
		-- Covering index for page view queries (includes all necessary columns)
		CREATE INDEX IF NOT EXISTS idx_pageview_covering ON events(timestamp DESC, event_name, session_id, is_bot);
		
		-- Index for session duration analytics
		CREATE INDEX IF NOT EXISTS idx_session_duration ON events(session_duration);`,
		Down: `DROP INDEX IF EXISTS idx_timestamp_is_bot;
		DROP INDEX IF EXISTS idx_timestamp_browser;
		DROP INDEX IF EXISTS idx_timestamp_device;
		DROP INDEX IF EXISTS idx_timestamp_os;
		DROP INDEX IF EXISTS idx_timestamp_referrer;
		DROP INDEX IF EXISTS idx_pageview_covering;
		DROP INDEX IF EXISTS idx_session_duration`,
	},
	{
		Version:     8,
		Description: "Add channel column to events table",
		Up:          `ALTER TABLE events ADD COLUMN IF NOT EXISTS channel VARCHAR`,
		Down:        `ALTER TABLE events DROP COLUMN IF EXISTS channel`,
	},
	{
		Version:     9,
		Description: "Create index on channel column",
		Up:          `CREATE INDEX IF NOT EXISTS idx_channel ON events(channel)`,
		Down:        `DROP INDEX IF EXISTS idx_channel`,
	},
	{
		Version:     10,
		Description: "Create composite index for channel analytics",
		Up:          `CREATE INDEX IF NOT EXISTS idx_timestamp_channel ON events(timestamp DESC, channel)`,
		Down:        `DROP INDEX IF EXISTS idx_timestamp_channel`,
	},
}

func initMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			description VARCHAR NOT NULL,
			applied_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
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
