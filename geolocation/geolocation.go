package geolocation

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

const (
	geoDBDir      = "data/geodb"
	geoDBFilename = "dbip-country.mmdb"
	geoDBPath     = geoDBDir + "/" + geoDBFilename
)

// GeoLocation represents geographic location data
type GeoLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
}

// Service handles IP geolocation lookups
type Service struct {
	db *maxminddb.Reader
}

// NewService creates a new geolocation service
func NewService() (*Service, error) {
	// Ensure database exists
	if err := ensureDatabase(); err != nil {
		return nil, fmt.Errorf("failed to ensure geolocation database: %w", err)
	}

	// Open the database
	db, err := maxminddb.Open(geoDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open geolocation database: %w", err)
	}

	log.Println("✓ Geolocation database loaded successfully")
	return &Service{db: db}, nil
}

// Close closes the geolocation database
func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Lookup performs IP geolocation lookup
func (s *Service) Lookup(ipStr string) (*GeoLocation, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	var record struct {
		Country struct {
			ISOCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
	}

	err = s.db.Lookup(ip).Decode(&record)
	if err != nil {
		return nil, fmt.Errorf("geolocation lookup failed: %w", err)
	}

	geo := &GeoLocation{
		CountryCode: record.Country.ISOCode,
		Country:     record.Country.Names["en"],
		City:        record.City.Names["en"],
	}

	if geo.Country == "Israel" {
		geo.Country = "Palestine"
		geo.CountryCode = "PS"
	}

	return geo, nil
}

// ensureDatabase checks if the database exists, downloads if missing
func ensureDatabase() error {
	// Check if database already exists
	if _, err := os.Stat(geoDBPath); err == nil {
		log.Println("✓ Geolocation database found")
		return nil
	}

	log.Println("⬇️  Geolocation database not found, downloading...")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(geoDBDir, 0755); err != nil {
		return fmt.Errorf("failed to create geodb directory: %w", err)
	}

	// Try to download the database
	if err := downloadDatabase(); err != nil {
		return fmt.Errorf("failed to download geolocation database: %w", err)
	}

	log.Println("✓ Geolocation database downloaded successfully")
	return nil
}

// downloadDatabase downloads the DB-IP Country Lite database
func downloadDatabase() error {
	// Try current month first, then last month
	now := time.Now()
	thisMonth := now.Format("2006-01")
	lastMonth := now.AddDate(0, -1, 0).Format("2006-01")

	urls := []string{
		fmt.Sprintf("https://download.db-ip.com/free/dbip-country-lite-%s.mmdb.gz", thisMonth),
		fmt.Sprintf("https://download.db-ip.com/free/dbip-country-lite-%s.mmdb.gz", lastMonth),
	}

	var lastErr error
	for i, url := range urls {
		log.Printf("Trying to download from: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Warning: failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode == 404 && i < len(urls)-1 {
			log.Printf("Got 404 for %s, trying next URL...", url)
			continue
		}

		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			continue
		}

		// Decompress and save
		if err := decompressAndSave(resp.Body); err != nil {
			return fmt.Errorf("failed to decompress database: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to download from all URLs: %w", lastErr)
}

// decompressAndSave decompresses the gzipped database and saves it
func decompressAndSave(body io.Reader) error {
	gzReader, err := gzip.NewReader(body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		if err := gzReader.Close(); err != nil {
			log.Printf("Warning: failed to close gzip reader: %v", err)
		}
	}()

	// Create temporary file
	tmpFile := geoDBPath + ".tmp"
	outFile, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			log.Printf("Warning: failed to close output file: %v", err)
		}
	}()

	// Copy decompressed data
	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		if removeErr := os.Remove(tmpFile); removeErr != nil {
			log.Printf("Warning: failed to remove temp file: %v", removeErr)
		}
		return fmt.Errorf("failed to write database: %w", err)
	}

	// Rename temp file to final name
	if err := os.Rename(tmpFile, geoDBPath); err != nil {
		if removeErr := os.Remove(tmpFile); removeErr != nil {
			log.Printf("Warning: failed to remove temp file: %v", removeErr)
		}
		return fmt.Errorf("failed to rename database file: %w", err)
	}

	return nil
}

// LookupOrDefault performs lookup and returns default values on error
func (s *Service) LookupOrDefault(ipStr string) *GeoLocation {
	geo, err := s.Lookup(ipStr)
	if err != nil {
		// Return unknown location on error
		return &GeoLocation{
			Country:     "Unknown",
			CountryCode: "XX",
			City:        "",
		}
	}
	return geo
}

// GetDatabasePath returns the path to the geolocation database
func GetDatabasePath() string {
	return filepath.Clean(geoDBPath)
}

// DatabaseExists checks if the geolocation database exists
func DatabaseExists() bool {
	_, err := os.Stat(geoDBPath)
	return err == nil
}
