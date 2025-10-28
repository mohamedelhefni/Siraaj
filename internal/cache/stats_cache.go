package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
	Size       int
	AccessTime int64 // For LRU eviction
}

type StatsCache struct {
	mu          sync.RWMutex
	items       map[string]*CacheItem
	maxSize     int64 // Max memory in bytes
	currentSize int64
	hits        int64
	misses      int64
}

// NewStatsCache creates a new cache with specified max size in MB
func NewStatsCache(maxSizeMB int) *StatsCache {
	c := &StatsCache{
		items:   make(map[string]*CacheItem),
		maxSize: int64(maxSizeMB) * 1024 * 1024,
	}

	// Start cleanup goroutine
	go c.cleanupExpired()

	// Log stats every hour
	go c.logStats()

	return c
}

// Set adds or updates a cache item
func (c *StatsCache) Set(key string, value interface{}, duration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Estimate size by marshaling to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	size := len(data)

	// Check if we need to evict items
	if c.currentSize+int64(size) > c.maxSize {
		c.evictLRU(int64(size))
	}

	// Remove old item size if updating existing key
	if oldItem, exists := c.items[key]; exists {
		c.currentSize -= int64(oldItem.Size)
	}

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
		Size:       size,
		AccessTime: time.Now().UnixNano(),
	}

	c.currentSize += int64(size)
	return nil
}

// Get retrieves a cache item if it exists and is not expired
func (c *StatsCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		c.misses++
		return nil, false
	}

	// Check expiration
	if time.Now().UnixNano() > item.Expiration {
		c.currentSize -= int64(item.Size)
		delete(c.items, key)
		c.misses++
		return nil, false
	}

	// Update access time for LRU
	item.AccessTime = time.Now().UnixNano()
	c.hits++
	return item.Value, true
}

// evictLRU removes least recently used items until we have enough space
func (c *StatsCache) evictLRU(neededSize int64) {
	type itemWithKey struct {
		key  string
		item *CacheItem
	}

	// Collect all items
	items := make([]itemWithKey, 0, len(c.items))
	for k, v := range c.items {
		items = append(items, itemWithKey{k, v})
	}

	// Sort by access time (oldest first) using bubble sort
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].item.AccessTime > items[j].item.AccessTime {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// Remove oldest items until we have space
	for _, item := range items {
		if c.currentSize+neededSize <= c.maxSize {
			break
		}
		c.currentSize -= int64(item.item.Size)
		delete(c.items, item.key)
	}
}

// cleanupExpired removes expired items periodically
func (c *StatsCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		expiredCount := 0
		for key, item := range c.items {
			if now > item.Expiration {
				c.currentSize -= int64(item.Size)
				delete(c.items, key)
				expiredCount++
			}
		}
		if expiredCount > 0 {
			log.Printf("[Cache] Cleaned up %d expired items\n", expiredCount)
		}
		c.mu.Unlock()
	}
}

// logStats logs cache statistics periodically
func (c *StatsCache) logStats() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.RLock()
		totalRequests := c.hits + c.misses
		hitRate := 0.0
		if totalRequests > 0 {
			hitRate = float64(c.hits) / float64(totalRequests) * 100
		}
		log.Printf("[Cache Stats] Items: %d, Size: %.2fMB/%.2fMB, Hit Rate: %.1f%%, Hits: %d, Misses: %d\n",
			len(c.items),
			float64(c.currentSize)/(1024*1024),
			float64(c.maxSize)/(1024*1024),
			hitRate,
			c.hits,
			c.misses)
		c.mu.RUnlock()
	}
}

// Clear removes all items from cache
func (c *StatsCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
	c.currentSize = 0
	c.hits = 0
	c.misses = 0
}

// GetStats returns current cache statistics
func (c *StatsCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalRequests := c.hits + c.misses
	hitRate := 0.0
	if totalRequests > 0 {
		hitRate = float64(c.hits) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"items":          len(c.items),
		"size_mb":        float64(c.currentSize) / (1024 * 1024),
		"max_size_mb":    float64(c.maxSize) / (1024 * 1024),
		"hit_rate":       hitRate,
		"hits":           c.hits,
		"misses":         c.misses,
		"total_requests": totalRequests,
	}
}

// GenerateKey creates a cache key from parameters
func GenerateKey(endpoint string, startDate, endDate time.Time, filters map[string]string, limit int) string {
	// Normalize timestamps to minute precision to improve cache hits
	start := startDate.Truncate(time.Minute).Unix()
	end := endDate.Truncate(time.Minute).Unix()

	data := fmt.Sprintf("%s|%d|%d|%v|%d", endpoint, start, end, filters, limit)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:16])
}
