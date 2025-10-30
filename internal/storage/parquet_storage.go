package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mohamedelhefni/siraaj/internal/domain"
)

const (
	// Default buffer size before flushing to disk
	DefaultBufferSize = 10000
	// Default flush interval
	DefaultFlushInterval = 30 * time.Second
	// Parquet file path
	DefaultParquetFile = "data/events.parquet"
	// Temp CSV path for buffering
	TempCSVFile = "data/events_buffer.csv"
)

// ParquetStorage handles buffered writes to a single Parquet file using DuckDB COPY
type ParquetStorage struct {
	db            *sql.DB
	filePath      string
	tempCSVPath   string
	buffer        []domain.Event
	bufferSize    int
	flushInterval time.Duration
	mu            sync.Mutex
	stopChan      chan struct{}
	flushChan     chan struct{}
	wg            sync.WaitGroup
	idCounter     uint64
}

// NewParquetStorage creates a new Parquet storage with buffering
func NewParquetStorage(db *sql.DB, filePath string, bufferSize int, flushInterval time.Duration) (*ParquetStorage, error) {
	if filePath == "" {
		filePath = DefaultParquetFile
	}
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}
	if flushInterval <= 0 {
		flushInterval = DefaultFlushInterval
	}

	// Ensure directory exists
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	ps := &ParquetStorage{
		db:            db,
		filePath:      filePath,
		tempCSVPath:   TempCSVFile,
		buffer:        make([]domain.Event, 0, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		stopChan:      make(chan struct{}),
		flushChan:     make(chan struct{}, 1),
		idCounter:     1,
	}

	// Start background flusher
	ps.wg.Add(1)
	go ps.backgroundFlusher()

	log.Printf("✓ Parquet storage initialized: file=%s, buffer_size=%d, flush_interval=%v",
		filePath, bufferSize, flushInterval)

	return ps, nil
}

// GetNextID returns the next ID for event insertion
func (ps *ParquetStorage) GetNextID() uint64 {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	id := ps.idCounter
	ps.idCounter++
	return id
}

// Write adds an event to the buffer
func (ps *ParquetStorage) Write(event domain.Event) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.buffer = append(ps.buffer, event)

	// Check if buffer is full
	if len(ps.buffer) >= ps.bufferSize {
		log.Printf("📦 Buffer full (%d events), triggering flush...", len(ps.buffer))
		// Trigger flush without blocking
		select {
		case ps.flushChan <- struct{}{}:
		default:
			// Flush already pending
		}
	}

	return nil
}

// WriteBatch adds multiple events to the buffer
func (ps *ParquetStorage) WriteBatch(events []domain.Event) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.buffer = append(ps.buffer, events...)

	// Check if buffer is full
	if len(ps.buffer) >= ps.bufferSize {
		log.Printf("📦 Buffer full (%d events), triggering flush...", len(ps.buffer))
		// Trigger flush without blocking
		select {
		case ps.flushChan <- struct{}{}:
		default:
			// Flush already pending
		}
	}

	return nil
}

// backgroundFlusher runs in a goroutine and flushes buffer periodically
func (ps *ParquetStorage) backgroundFlusher() {
	defer ps.wg.Done()

	ticker := time.NewTicker(ps.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ps.stopChan:
			// Final flush before shutdown
			if err := ps.Flush(); err != nil {
				log.Printf("❌ Error during final flush: %v", err)
			}
			return

		case <-ticker.C:
			// Periodic flush
			if err := ps.Flush(); err != nil {
				log.Printf("❌ Error during periodic flush: %v", err)
			}

		case <-ps.flushChan:
			// Manual flush triggered by full buffer
			if err := ps.Flush(); err != nil {
				log.Printf("❌ Error during manual flush: %v", err)
			}
		}
	}
}

// Flush writes buffered events to Parquet file using DuckDB COPY
func (ps *ParquetStorage) Flush() error {
	ps.mu.Lock()
	if len(ps.buffer) == 0 {
		ps.mu.Unlock()
		return nil
	}

	// Copy buffer and clear it
	eventsToWrite := make([]domain.Event, len(ps.buffer))
	copy(eventsToWrite, ps.buffer)
	ps.buffer = ps.buffer[:0]
	ps.mu.Unlock()

	start := time.Now()
	log.Printf("💾 Flushing %d events to Parquet file...", len(eventsToWrite))

	// Write events to temporary CSV file
	csvFile, err := os.Create(ps.tempCSVPath)
	if err != nil {
		return fmt.Errorf("failed to create temp CSV: %w", err)
	}
	defer func() {
		csvFile.Close()
		os.Remove(ps.tempCSVPath) // Clean up temp file
	}()

	// Write CSV data
	fmt.Fprintf(csvFile, "id,timestamp,event_name,user_id,session_id,session_duration,url,referrer,user_agent,ip,country,browser,os,device,is_bot,project_id,channel\n")
	for _, event := range eventsToWrite {
		// Format timestamp as ISO8601 string for DuckDB
		timestampStr := event.Timestamp.UTC().Format("2006-01-02 15:04:05.000000")
		fmt.Fprintf(csvFile, "%d,%s,%s,%s,%s,%d,%s,%s,%s,%s,%s,%s,%s,%s,%t,%s,%s\n",
			event.ID,
			timestampStr,
			escapeCsv(event.EventName),
			escapeCsv(event.UserID),
			escapeCsv(event.SessionID),
			event.SessionDuration,
			escapeCsv(event.URL),
			escapeCsv(event.Referrer),
			escapeCsv(event.UserAgent),
			escapeCsv(event.IP),
			escapeCsv(event.Country),
			escapeCsv(event.Browser),
			escapeCsv(event.OS),
			escapeCsv(event.Device),
			event.IsBot,
			escapeCsv(event.ProjectID),
			escapeCsv(event.Channel),
		)
	}
	csvFile.Close()

	// Use DuckDB COPY to convert CSV to Parquet with ZSTD compression
	// Sort by timestamp for better query performance (row group pruning)
	var copyQuery string
	if _, err := os.Stat(ps.filePath); os.IsNotExist(err) {
		// File doesn't exist, create new (sorted by timestamp)
		copyQuery = fmt.Sprintf(`
			COPY (
				SELECT * FROM read_csv('%s', 
					AUTO_DETECT=TRUE,
					header=true,
					timestampformat='%%Y-%%m-%%d %%H:%%M:%%S.%%f'
				)
				ORDER BY timestamp
			) TO '%s' (FORMAT 'PARQUET', CODEC 'ZSTD', ROW_GROUP_SIZE 100000)
		`, ps.tempCSVPath, ps.filePath)
	} else {
		// File exists, append to it (keep sorted by timestamp)
		copyQuery = fmt.Sprintf(`
			COPY (
				SELECT * FROM (
					SELECT * FROM read_parquet('%s')
					UNION ALL
					SELECT * FROM read_csv('%s', 
						AUTO_DETECT=TRUE,
						header=true,
						timestampformat='%%Y-%%m-%%d %%H:%%M:%%S.%%f'
					)
				)
				ORDER BY timestamp
			) TO '%s' (FORMAT 'PARQUET', CODEC 'ZSTD', ROW_GROUP_SIZE 100000)
		`, ps.filePath, ps.tempCSVPath, ps.filePath)

		// Backup and remove old file
		backupPath := ps.filePath + ".backup"
		os.Rename(ps.filePath, backupPath)
		defer os.Remove(backupPath)
	}

	_, err = ps.db.Exec(copyQuery)
	if err != nil {
		return fmt.Errorf("failed to copy to Parquet: %w", err)
	}

	duration := time.Since(start)
	log.Printf("✅ Flushed %d events to Parquet in %v (%.0f events/sec)",
		len(eventsToWrite), duration, float64(len(eventsToWrite))/duration.Seconds())

	return nil
}

// escapeCsv escapes CSV fields
func escapeCsv(s string) string {
	if s == "" {
		return ""
	}
	// Simple CSV escaping - quote fields with commas or quotes
	needsQuote := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuote = true
			break
		}
	}
	if needsQuote {
		// Escape quotes by doubling them
		escaped := ""
		for _, c := range s {
			if c == '"' {
				escaped += "\"\""
			} else {
				escaped += string(c)
			}
		}
		return "\"" + escaped + "\""
	}
	return s
}

// Close gracefully shuts down the storage, flushing any remaining data
func (ps *ParquetStorage) Close() error {
	log.Println("🛑 Shutting down Parquet storage...")

	// Stop background flusher
	close(ps.stopChan)

	// Wait for background flusher to complete
	ps.wg.Wait()

	log.Println("✓ Parquet storage shut down successfully")
	return nil
}

// GetFilePath returns the Parquet file path for DuckDB queries
func (ps *ParquetStorage) GetFilePath() string {
	return ps.filePath
}
