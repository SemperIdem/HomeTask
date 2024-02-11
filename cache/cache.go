package cache

import (
	"bytes"
	"homeTask/config"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Entry struct {
	entry []byte
	ttl   time.Time
}

type MemoryCache struct {
	store map[string]Entry
	sync.RWMutex
}

func New() *MemoryCache {
	return &MemoryCache{
		store: map[string]Entry{},
	}
}

func (m *MemoryCache) Get(key string) ([]byte, bool) {
	m.RLock()
	defer m.RUnlock()
	entry, ok := m.store[key]
	if !ok || time.Now().After(entry.ttl) {
		return nil, false
	}
	return entry.entry, ok
}

func (m *MemoryCache) Set(key string, entry []byte, ttl time.Time) {
	m.Lock()
	defer m.Unlock()
	m.store[key] = Entry{
		entry: entry,
		ttl:   ttl,
	}
}

type CachedTransport struct {
	Cache         *MemoryCache
	BaseTransport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction, caching the response body.
func (t *CachedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Generate a cache key based on the request
	cacheKey := req.URL.String()

	// Check if the response body is already cached
	if cachedBody, ok := t.Cache.Get(cacheKey); ok {
		log.Println("Cache hit for", cacheKey)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(cachedBody)),
		}, nil
	}

	// If not cached, make the request and cache the response body
	log.Println("Cache miss for", cacheKey)
	resp, err := t.BaseTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the response is empty or an error status code
	if resp.StatusCode >= 400 || len(body) == 0 {
		return resp, nil // Don't cache error responses or empty bodies
	}

	// Cache the response body
	t.Cache.Set(cacheKey, body, time.Now().Add(config.DefaultCacheTTL)) // Adjust expiry time as needed

	// Since we've read the response body, we need to create a new response
	// with the cached body and the original status code
	return &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}, nil
}
