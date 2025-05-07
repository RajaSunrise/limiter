package limiter

import (
	"context"
	"sync"
	"time"
)

// Store defines the interface for limiter
// Store have 4 values Take, Rollback, Get and Set
type Store interface {
	Take(ctx context.Context, key string, maxRequests int, window time.Duration, algorithm string) (bool, int, time.Time, error)
	Rollback(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, value int, expiration time.Duration) error
}

type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]*MemoryEntries
}

type MemoryEntries struct {
	count     int
	expiresAt time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make(map[string]*MemoryEntries),
	}
}

func (m *MemoryStore) Take(ctx context.Context, key string, maxRequests int, window time.Duration, algorithm string) (bool, int, time.Time, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	reset := now.Add(window)

	// Cleanup expired entries
	for k, v := range m.entries {
		if now.After(v.expiresAt) {
			delete(m.entries, k)
		}
	}

	entry, exists := m.entries[key]
	if !exists {
		m.entries[key] = &MemoryEntries{
			count:     1,
			expiresAt: reset,
		}
		return true, maxRequests - 1, reset, nil
	}

	if entry.count >= maxRequests {
		return false, 0, entry.expiresAt, nil
	}

	entry.count++
	return true, maxRequests - entry.count, entry.expiresAt, nil
}

func (m *MemoryStore) Rollback(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, exists := m.entries[key]; exists {
		entry.count--
		if entry.count <= 0 {
			delete(m.entries, key)
		}
	}
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, exists := m.entries[key]; exists {
		return entry.count, nil
	}
	return 0, nil
}

func (m *MemoryStore) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries[key] = &MemoryEntries{
		count:     value,
		expiresAt: time.Now().Add(expiration),
	}
	return nil
}

func (m *MemoryStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = make(map[string]*MemoryEntries)
	return nil
}
