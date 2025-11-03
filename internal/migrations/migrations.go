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
		Description: "Create events table with optimized schema",
		Up: `CREATE TABLE IF NOT EXISTS events (
			id UBIGINT PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			date_hour TIMESTAMP NOT NULL,
			date_day DATE NOT NULL,
			date_month DATE NOT NULL,
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
			is_bot BOOLEAN DEFAULT FALSE,
			project_id VARCHAR DEFAULT 'default',
			channel VARCHAR
		);
		CREATE SEQUENCE IF NOT EXISTS id_sequence START 1;`,
		Down: `DROP SEQUENCE IF EXISTS id_sequence;
		DROP TABLE IF EXISTS events;`,
	},
	{
		Version:     2,
		Description: "Create id sequence",
		Up:          `CREATE SEQUENCE IF NOT EXISTS id_sequence START 1`,
		Down:        `DROP SEQUENCE IF EXISTS id_sequence`,
	},
	{
		Version:     3,
		Description: "Create optimized indexes for large-scale analytics",
		Up: `-- Primary time-based indexes using date partitioning columns
		CREATE INDEX IF NOT EXISTS idx_date_day ON events(date_day DESC);
		CREATE INDEX IF NOT EXISTS idx_date_hour ON events(date_hour DESC);
		
		-- Covering indexes for common analytics queries
		CREATE INDEX IF NOT EXISTS idx_day_project_event ON events(date_day, project_id, event_name, is_bot);
		CREATE INDEX IF NOT EXISTS idx_day_country ON events(date_day, country);
		CREATE INDEX IF NOT EXISTS idx_day_channel ON events(date_day, channel);
		
		-- Session and user analysis indexes
		CREATE INDEX IF NOT EXISTS idx_session_timestamp ON events(session_id, timestamp) ;
		CREATE INDEX IF NOT EXISTS idx_user_day ON events(user_id, date_day) ;
		
		-- URL and referrer analysis indexes
		CREATE INDEX IF NOT EXISTS idx_day_url ON events(date_day, url);
		CREATE INDEX IF NOT EXISTS idx_day_referrer ON events(date_day, referrer);
		
		-- Device analytics indexes
		CREATE INDEX IF NOT EXISTS idx_day_browser ON events(date_day, browser);
		CREATE INDEX IF NOT EXISTS idx_day_device ON events(date_day, device) ;
		CREATE INDEX IF NOT EXISTS idx_day_os ON events(date_day, os) ;`,
		Down: `DROP INDEX IF EXISTS idx_date_day;
		DROP INDEX IF EXISTS idx_date_hour;
		DROP INDEX IF EXISTS idx_day_project_event;
		DROP INDEX IF EXISTS idx_day_country;
		DROP INDEX IF EXISTS idx_day_channel;
		DROP INDEX IF EXISTS idx_session_timestamp;
		DROP INDEX IF EXISTS idx_user_day;
		DROP INDEX IF EXISTS idx_day_url;
		DROP INDEX IF EXISTS idx_day_referrer;
		DROP INDEX IF EXISTS idx_day_browser;
		DROP INDEX IF EXISTS idx_day_device;
		DROP INDEX IF EXISTS idx_day_os`,
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
