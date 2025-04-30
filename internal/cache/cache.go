package cache

import (
	"sync"
	"time"
	"github.com/yourusername/exchange-rate-service/internal/model"
)

// RateCache is a thread-safe in-memory cache for exchange rates
type RateCache struct {
	mu    sync.RWMutex
	cache map[string]model.ExchangeRate // key: base_target_date
}

func NewRateCache() *RateCache {
	return &RateCache{
		cache: make(map[string]model.ExchangeRate),
	}
}

func cacheKey(base, target, date string) string {
	return base + "_" + target + "_" + date
}

func (rc *RateCache) Get(base, target, date string) (model.ExchangeRate, bool) {
	key := cacheKey(base, target, date)
	rc.mu.RLock()
	rate, ok := rc.cache[key]
	rc.mu.RUnlock()
	return rate, ok
}

func (rc *RateCache) Set(rate model.ExchangeRate) {
	key := cacheKey(rate.Base, rate.Target, rate.Date)
	rc.mu.Lock()
	rc.cache[key] = rate
	rc.mu.Unlock()
}

// Optional: Purge old entries (not strictly required for this assignment)
func (rc *RateCache) PurgeOlderThan(d time.Duration) {
	cutoff := time.Now().Add(-d)
	rc.mu.Lock()
	for k, v := range rc.cache {
		if v.UpdatedAt.Before(cutoff) {
			delete(rc.cache, k)
		}
	}
	rc.mu.Unlock()
}
