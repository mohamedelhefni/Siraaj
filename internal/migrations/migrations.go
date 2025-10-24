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
			url VARCHAR,
			referrer VARCHAR,
			user_agent VARCHAR,
			ip VARCHAR,
			country VARCHAR,
			browser VARCHAR,
			os VARCHAR,
			device VARCHAR,
			properties JSON,
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
		Up: `CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp);
		CREATE INDEX IF NOT EXISTS idx_event_name ON events(event_name);
		CREATE INDEX IF NOT EXISTS idx_user_id ON events(user_id);
		CREATE INDEX IF NOT EXISTS idx_country ON events(country);
		CREATE INDEX IF NOT EXISTS idx_referrer ON events(referrer);
		CREATE INDEX IF NOT EXISTS idx_project_id ON events(project_id)`,
		Down: `DROP INDEX IF EXISTS idx_timestamp;
		DROP INDEX IF EXISTS idx_event_name;
		DROP INDEX IF EXISTS idx_user_id;
		DROP INDEX IF EXISTS idx_country;
		DROP INDEX IF EXISTS idx_referrer;
		DROP INDEX IF EXISTS idx_project_id`,
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
