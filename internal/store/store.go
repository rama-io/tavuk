package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type Target struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	RoleID    string `json:"role_id"`
}

// SON-file-backed mapping
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Target
}

// this loads the store from path, starting empty if the file is missing.
func Open(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: map[string]Target{},
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read store: %w", err)
	}
	if len(raw) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(raw, &s.data); err != nil {
		return nil, fmt.Errorf("parse store: %w", err)
	}
	return s, nil
}

// Get returns the target for an app and whether it exists.
func (s *Store) Get(app string) (Target, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.data[app]
	return t, ok
}

// Apps returns the stored app names in sorted order.
func (s *Store) Apps() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.data))
	for name := range s.data {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Set stores an app's target and persists the change atomically.
func (s *Store) Set(app string, t Target) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[app] = t

	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode store: %w", err)
	}

	// Sir we need a backup!
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return fmt.Errorf("write store: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("atomic replace store: %w", err)
	}
	return nil
}
