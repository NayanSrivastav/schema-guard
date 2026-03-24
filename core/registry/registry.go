package registry

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrSchemaNotFound is returned when querying an unknown schema by name or version
	ErrSchemaNotFound = errors.New("schema not found")
	// ErrVersionExists ensures we don't destructively overwrite locked schema versions
	ErrVersionExists  = errors.New("schema version already exists")
)

// SchemaRecord represents a distinct version-controlled schema definition
type SchemaRecord struct {
	Name        string
	Version     string // e.g., "v1.0", "v2.0"
	SchemaJSON  string
	CreatedAt   time.Time
	Description string
}

// Registry defines the storage and retrieval interface for schema catalogs
type Registry interface {
	Save(SchemaRecord) error
	Get(name, version string) (*SchemaRecord, error)
	GetLatest(name string) (*SchemaRecord, error)
	ListVersions(name string) ([]string, error)
}

// MemoryRegistry is a fast, thread-safe, in-memory capability for MVP testing
type MemoryRegistry struct {
	mu      sync.RWMutex
	schemas map[string]map[string]SchemaRecord // map[name]map[version]SchemaRecord
}

// NewMemoryRegistry constructs a new instance of the InMemory cache
func NewMemoryRegistry() *MemoryRegistry {
	return &MemoryRegistry{
		schemas: make(map[string]map[string]SchemaRecord),
	}
}

// Save locks the specific snapshot into the registry
func (m *MemoryRegistry) Save(rec SchemaRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.schemas[rec.Name]; !exists {
		m.schemas[rec.Name] = make(map[string]SchemaRecord)
	}

	if _, exists := m.schemas[rec.Name][rec.Version]; exists {
		return ErrVersionExists
	}

	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = time.Now()
	}

	m.schemas[rec.Name][rec.Version] = rec
	return nil
}

// Get selects the exactly matching named schema and semantic version
func (m *MemoryRegistry) Get(name, version string) (*SchemaRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, exists := m.schemas[name]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	rec, exists := versions[version]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	return &rec, nil
}

// GetLatest acts as a chronologically-ordered fallback for un-pinned models
func (m *MemoryRegistry) GetLatest(name string) (*SchemaRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, exists := m.schemas[name]
	if !exists || len(versions) == 0 {
		return nil, ErrSchemaNotFound
	}

	var latest SchemaRecord
	var latestTime time.Time
	for _, rec := range versions {
		if rec.CreatedAt.After(latestTime) || latestTime.IsZero() {
			latest = rec
			latestTime = rec.CreatedAt
		}
	}

	return &latest, nil
}

// ListVersions yields all active iterations available for a given schema
func (m *MemoryRegistry) ListVersions(name string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, exists := m.schemas[name]
	if !exists {
		return nil, ErrSchemaNotFound
	}

	var list []string
	for v := range versions {
		list = append(list, v)
	}
	return list, nil
}
