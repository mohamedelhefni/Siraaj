package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/mohamedelhefni/siraaj/geolocation"
	"github.com/mohamedelhefni/siraaj/internal/handler"
	"github.com/mohamedelhefni/siraaj/internal/middleware"
	"github.com/mohamedelhefni/siraaj/internal/migrations"
	"github.com/mohamedelhefni/siraaj/internal/repository"
	"github.com/mohamedelhefni/siraaj/internal/service"
)

//go:embed all:ui/dashboard
var uiFiles embed.FS

// initDatabase initializes the database connection and runs migrations
func initDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	// Enable DuckDB optimizations
	if _, err = db.Exec("PRAGMA memory_limit='512MB'"); err != nil {
		log.Printf("Warning: Could not set memory limit: %v", err)
	}

	if _, err = db.Exec("PRAGMA threads=2"); err != nil {
		log.Printf("Warning: Could not set threads: %v", err)
	}

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return db, nil
}

// getClientIP extracts the real client IP from request
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
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
		defer geoService.Close()
	}

	// Initialize database
	db, err := initDatabase("analytics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("✓ DuckDB initialized successfully")

	// Initialize layers
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService, geoService)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/track", eventHandler.TrackEvent)
	mux.HandleFunc("/api/stats", eventHandler.GetStats)
	mux.HandleFunc("/api/events", eventHandler.GetEvents)
	mux.HandleFunc("/api/online", eventHandler.GetOnlineUsers)
	mux.HandleFunc("/api/projects", eventHandler.GetProjects)
	mux.HandleFunc("/api/health", eventHandler.Health)
	mux.HandleFunc("/api/geo", eventHandler.GeoTest)

	// Debug endpoint to show all events
	mux.HandleFunc("/api/debug/events", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, timestamp, event_name, user_id FROM events ORDER BY timestamp DESC LIMIT 50")
		if err != nil {
			log.Printf("Error querying events: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"events": events,
			"count":  len(events),
		})
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

	port := "8080"
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 Analytics Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🎨 Dashboard:  http://localhost:%s/dashboard/\n", port)
	fmt.Printf("📡 API Track:  http://localhost:%s/api/track\n", port)
	fmt.Printf("📈 API Stats:  http://localhost:%s/api/stats\n", port)
	fmt.Printf("🌍 Geo Test:   http://localhost:%s/api/geo\n", port)
	fmt.Printf("❤️  Health:    http://localhost:%s/api/health\n", port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Server ready - Using official DuckDB Go driver")
	fmt.Println("✓ Svelte Dashboard embedded and ready")
	if geoService != nil {
		fmt.Println("✓ Geolocation service enabled")
	} else {
		fmt.Println("⚠️  Geolocation service disabled")
	}
	fmt.Println("✓ Clean Architecture implemented")
	fmt.Println()

	// Apply middleware: CORS and Logging
	httpHandler := middleware.CORS(middleware.Logging(mux))
	log.Fatal(http.ListenAndServe(":"+port, httpHandler))
}
