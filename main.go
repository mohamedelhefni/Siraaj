package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/mohamedelhefni/siraaj/geolocation"
	"github.com/mohamedelhefni/siraaj/internal/handler"
	"github.com/mohamedelhefni/siraaj/internal/middleware"
	"github.com/mohamedelhefni/siraaj/internal/migrations"
	"github.com/mohamedelhefni/siraaj/internal/repository"
	"github.com/mohamedelhefni/siraaj/internal/service"
	"github.com/mohamedelhefni/siraaj/internal/storage"
)

//go:embed all:ui/dashboard
var uiFiles embed.FS

// initDatabase initializes the database connection and runs migrations
func initDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// ClickHouse settings are configured via DSN or SET queries
	// Example: Set max_memory_usage if needed
	maxMemory := os.Getenv("CLICKHOUSE_MAX_MEMORY")
	if maxMemory != "" {
		if _, err = db.Exec(fmt.Sprintf("SET max_memory_usage = %s", maxMemory)); err != nil {
			log.Printf("Warning: Could not set max_memory_usage: %v", err)
		} else {
			log.Printf("✓ ClickHouse max_memory_usage set to: %s", maxMemory)
		}
	}

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return db, nil
}

func main() {
	// Initialize geolocation service
	geoService, err := geolocation.NewService()
	if err != nil {
		log.Printf("⚠️  Warning: Geolocation service unavailable: %v", err)
		log.Println("⚠️  Continuing without geolocation support...")
		geoService = nil
	}
	if geoService != nil {
		defer func() {
			if err := geoService.Close(); err != nil {
				log.Printf("Warning: failed to close geolocation service: %v", err)
			}
		}()
	}

	// Initialize database first (needed for Parquet storage)
	clickhouseDSN := os.Getenv("CLICKHOUSE_DSN")
	if clickhouseDSN == "" {
		clickhouseDSN = "clickhouse://localhost:9000/siraaj?username=default&password="
	}

	db, err := initDatabase(clickhouseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Warning: failed to close database: %v", err)
		}
	}()

	log.Println("✓ ClickHouse initialized successfully")

	// Initialize storage layer (no-op for ClickHouse, kept for compatibility)
	parquetFilePath := os.Getenv("PARQUET_FILE")
	if parquetFilePath == "" {
		parquetFilePath = "data/events"
	}

	bufferSize := 10000               // Not used with ClickHouse
	flushInterval := 30 * time.Second // Not used with ClickHouse

	parquetStorage, err := storage.NewParquetStorage(db, parquetFilePath, bufferSize, flushInterval)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer func() {
		if err := parquetStorage.Close(); err != nil {
			log.Printf("Warning: failed to close storage: %v", err)
		}
	}()

	// Initialize layers with Parquet storage and caching
	baseRepo := repository.NewEventRepository(db, parquetStorage)

	eventService := service.NewEventService(baseRepo)
	eventHandler := handler.NewEventHandler(eventService, geoService)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\n🛑 Shutting down gracefully...")

		// Close storage first
		if err := parquetStorage.Close(); err != nil {
			log.Printf("Error closing storage: %v", err)
		}

		// Close other resources
		if geoService != nil {
			if err := geoService.Close(); err != nil {
				log.Printf("Error closing geolocation service: %v", err)
			}
		}

		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}

		os.Exit(0)
	}()

	// Setup HTTP routes
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/track", eventHandler.TrackEvent)
	mux.HandleFunc("/api/track/batch", eventHandler.TrackBatchEvents)
	mux.HandleFunc("/api/stats", eventHandler.GetStats)
	mux.HandleFunc("/api/events", eventHandler.GetEvents)
	mux.HandleFunc("/api/online", eventHandler.GetOnlineUsers)
	mux.HandleFunc("/api/projects", eventHandler.GetProjects)
	mux.HandleFunc("/api/funnel", eventHandler.GetFunnelAnalysis)
	mux.HandleFunc("/api/health", eventHandler.Health)
	mux.HandleFunc("/api/geo", eventHandler.GeoTest)

	// New focused stats endpoints
	mux.HandleFunc("/api/stats/overview", eventHandler.GetTopStats)
	mux.HandleFunc("/api/stats/timeline", eventHandler.GetTimeline)
	mux.HandleFunc("/api/stats/pages", eventHandler.GetTopPagesHandler)
	mux.HandleFunc("/api/stats/pages/entry-exit", eventHandler.GetEntryExitPagesHandler)
	mux.HandleFunc("/api/stats/countries", eventHandler.GetTopCountriesHandler)
	mux.HandleFunc("/api/stats/sources", eventHandler.GetTopSourcesHandler)
	mux.HandleFunc("/api/stats/events", eventHandler.GetTopEventsHandler)
	mux.HandleFunc("/api/stats/devices", eventHandler.GetBrowsersDevicesOSHandler)

	// Channel analytics
	mux.HandleFunc("/api/channels", eventHandler.GetChannelsHandler)

	// Debug endpoint to show all events
	mux.HandleFunc("/api/debug/events", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, timestamp, event_name, user_id FROM events ORDER BY timestamp DESC LIMIT 50")
		if err != nil {
			log.Printf("Error querying events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				log.Printf("Warning: failed to close rows: %v", err)
			}
		}()

		events := []map[string]interface{}{}
		for rows.Next() {
			var id uint64
			var timestamp time.Time
			var eventName, userID string
			if err := rows.Scan(&id, &timestamp, &eventName, &userID); err != nil {
				continue
			}
			events = append(events, map[string]interface{}{
				"id":         id,
				"timestamp":  timestamp.Format(time.RFC3339),
				"event_name": eventName,
				"user_id":    userID,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"events": events,
			"count":  len(events),
		}); err != nil {
			log.Printf("Error encoding debug events: %v", err)
		}
	})

	// Storage stats endpoint
	mux.HandleFunc("/api/debug/storage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"storage_type": "ClickHouse native",
			"database":     "ClickHouse",
		}); err != nil {
			log.Printf("Error encoding storage stats: %v", err)
		}
	})

	// Serve dashboard (SvelteKit app)
	dashboardFS, err := fs.Sub(uiFiles, "ui/dashboard")
	if err != nil {
		log.Printf("Warning: Could not load dashboard: %v", err)
	} else {
		mux.Handle("/dashboard/", http.StripPrefix("/dashboard", http.FileServer(http.FS(dashboardFS))))
	}

	// Serve UI (must be last as it's a catch-all)
	mux.Handle("/", http.FileServer(http.FS(uiFiles)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Analytics Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🎨 Dashboard:  http://localhost:%s/dashboard/\n", port)
	fmt.Printf("📡 API Track:  http://localhost:%s/api/track\n", port)
	fmt.Printf("📦 API Batch:  http://localhost:%s/api/track/batch\n", port)
	fmt.Printf("📈 API Stats:  http://localhost:%s/api/stats\n", port)
	fmt.Printf("🌍 Geo Test:   http://localhost:%s/api/geo\n", port)
	fmt.Printf("❤️  Health:    http://localhost:%s/api/health\n", port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Server ready - Using ClickHouse database")
	fmt.Println("✓ Svelte Dashboard embedded and ready")
	if geoService != nil {
		fmt.Println("✓ Geolocation service enabled")
	} else {
		fmt.Println("⚠️  Geolocation service disabled")
	}
	fmt.Println("✓ Clean Architecture implemented")
	fmt.Println("✓ Using ClickHouse columnar database for high-performance analytics")
	fmt.Println()

	// Apply middleware: CORS and Logging
	httpHandler := middleware.CORS(middleware.Logging(mux))
	log.Fatal(http.ListenAndServe(":"+port, httpHandler))
}
