package nem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Cache stores downloaded report files on disk.
type Cache struct {
	dir string
}

// NewCache creates a cache rooted at dir, creating it when missing.
func NewCache(dir string) (*Cache, error) {
	if dir == "" {
		return nil, fmt.Errorf("nem: cache dir cannot be empty")
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("nem: create cache dir: %w", err)
	}
	return &Cache{dir: dir}, nil
}

// Path returns the on-disk path for a file reference.
func (c *Cache) Path(ref FileRef) string {
	return filepath.Join(c.dir, sanitize(ref.Report.Dir), sanitize(ref.Name))
}

// Get returns the cached bytes for ref, reporting whether it was present.
func (c *Cache) Get(ref FileRef) ([]byte, bool) {
	data, err := os.ReadFile(c.Path(ref))
	if err != nil {
		return nil, false
	}
	return data, true
}

// Put writes data for ref atomically.
func (c *Cache) Put(ref FileRef, data []byte) error {
	target := c.Path(ref)
	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return fmt.Errorf("nem: create cache subdir: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(target), ".tmp-*")
	if err != nil {
		return fmt.Errorf("nem: create temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("nem: write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("nem: close temp file: %w", err)
	}
	if err := os.Rename(tmpName, target); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("nem: rename temp file: %w", err)
	}
	return nil
}

func sanitize(name string) string {
	return strings.NewReplacer("/", "_", "\\", "_", "..", "_").Replace(name)
}
