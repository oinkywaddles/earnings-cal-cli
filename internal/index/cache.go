package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDir  = ".cache/earnings-cal-cli"
	cacheFile = "index_constituents.json"
	cacheTTL  = 90 * 24 * time.Hour
)

// Constituents holds index constituent symbols.
type Constituents struct {
	UpdatedAt time.Time `json:"updated_at"`
	SP100     []string  `json:"sp100"`
	Nasdaq100 []string  `json:"nasdaq100"`
	DowJones  []string  `json:"dowjones"`
	Merged    []string  `json:"merged"`
}

func cachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, cacheDir, cacheFile)
}

// GetConstituents returns constituents, using cache if available and fresh.
func GetConstituents(forceRefresh bool) (*Constituents, error) {
	if !forceRefresh {
		if cached, err := loadCache(); err == nil {
			return cached, nil
		}
	}

	c, err := FetchAllConstituents()
	if err != nil {
		return nil, err
	}

	_ = saveCache(c)
	return c, nil
}

// CacheStatus describes the state of the local cache.
type CacheStatus struct {
	Exists    bool
	Expired   bool
	UpdatedAt time.Time
	Path      string
	Count     int // number of merged symbols
}

// GetCacheStatus checks the current cache state without loading full data.
func GetCacheStatus() CacheStatus {
	path := cachePath()
	info, err := os.Stat(path)
	if err != nil {
		return CacheStatus{Path: path}
	}

	s := CacheStatus{
		Exists: true,
		Path:   path,
	}

	if time.Since(info.ModTime()) > cacheTTL {
		s.Expired = true
	}

	if data, err := os.ReadFile(path); err == nil {
		var c Constituents
		if json.Unmarshal(data, &c) == nil {
			s.UpdatedAt = c.UpdatedAt
			s.Count = len(c.Merged)
		}
	}

	return s
}

func loadCache() (*Constituents, error) {
	path := cachePath()

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if time.Since(info.ModTime()) > cacheTTL {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Constituents
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func saveCache(c *Constituents) error {
	path := cachePath()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
