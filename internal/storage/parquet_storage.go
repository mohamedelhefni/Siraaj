package storage

import (
	"database/sql"
	"log"
	"time"

	"github.com/mohamedelhefni/siraaj/internal/domain"
)

// ParquetStorage is now a no-op for ClickHouse since ClickHouse handles storage natively
// Kept for compatibility with existing code structure
type ParquetStorage struct {
	db      *sql.DB
	dataDir string
}

// NewParquetStorage creates a new storage instance (no-op for ClickHouse)
func NewParquetStorage(db *sql.DB, dataDir string, bufferSize int, flushInterval time.Duration) (*ParquetStorage, error) {
	log.Println("✓ Storage initialized (using ClickHouse native storage)")
	return &ParquetStorage{
		db:      db,
		dataDir: dataDir,
	}, nil
}

// Write is a no-op for ClickHouse (events are written directly to the database)
func (ps *ParquetStorage) Write(event domain.Event) error {
	return nil
}

// WriteBatch is a no-op for ClickHouse (events are written directly to the database)
func (ps *ParquetStorage) WriteBatch(events []domain.Event) error {
	return nil
}

// Flush is a no-op for ClickHouse
func (ps *ParquetStorage) Flush() error {
	return nil
}

// Close gracefully shuts down the storage
func (ps *ParquetStorage) Close() error {
	log.Println("✓ Storage closed")
	return nil
}

// GetFilePath returns an empty string (not used with ClickHouse)
func (ps *ParquetStorage) GetFilePath() string {
	return ""
}

// GetFileCount returns 0 (not applicable for ClickHouse)
func (ps *ParquetStorage) GetFileCount() (int, error) {
	return 0, nil
}

// GetNextID is not needed (handled by repository)
func (ps *ParquetStorage) GetNextID() uint64 {
	return 0
}
