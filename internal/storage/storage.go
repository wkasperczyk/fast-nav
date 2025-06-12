package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

type Bookmark struct {
	Path      string    `json:"path"`
	Created   time.Time `json:"created"`
	UsedCount int       `json:"used_count"`
	LastUsed  time.Time `json:"last_used"`
}

type BookmarkData struct {
	Version   string               `json:"version"`
	Bookmarks map[string]*Bookmark `json:"bookmarks"`
}

type Store struct {
	configDir string
	filePath  string
	data      *BookmarkData
}

func NewStore() (*Store, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".fn")
	filePath := filepath.Join(configDir, "bookmarks.json")
	
	store := &Store{
		configDir: configDir,
		filePath:  filePath,
	}
	
	err = store.ensureConfigDir()
	if err != nil {
		return nil, err
	}
	
	err = store.load()
	if err != nil {
		return nil, err
	}
	
	return store, nil
}

func (s *Store) ensureConfigDir() error {
	if _, err := os.Stat(s.configDir); os.IsNotExist(err) {
		return os.MkdirAll(s.configDir, 0755)
	}
	return nil
}

func (s *Store) load() error {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Create empty data structure
		s.data = &BookmarkData{
			Version:   "1.0",
			Bookmarks: make(map[string]*Bookmark),
		}
		return s.save()
	}
	
	file, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read bookmarks file: %w", err)
	}
	
	s.data = &BookmarkData{}
	err = json.Unmarshal(file, s.data)
	if err != nil {
		return fmt.Errorf("failed to parse bookmarks file: %w", err)
	}
	
	// Initialize bookmarks map if nil
	if s.data.Bookmarks == nil {
		s.data.Bookmarks = make(map[string]*Bookmark)
	}
	
	return nil
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bookmarks: %w", err)
	}
	
	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write bookmarks file: %w", err)
	}
	
	return nil
}

func (s *Store) SaveBookmark(alias, path string) error {
	now := time.Now()
	
	if existing, exists := s.data.Bookmarks[alias]; exists {
		// Update existing bookmark
		existing.Path = path
		existing.LastUsed = now
	} else {
		// Create new bookmark
		s.data.Bookmarks[alias] = &Bookmark{
			Path:      path,
			Created:   now,
			UsedCount: 0,
			LastUsed:  now,
		}
	}
	
	return s.save()
}

func (s *Store) GetBookmark(alias string) (*Bookmark, bool) {
	bookmark, exists := s.data.Bookmarks[alias]
	return bookmark, exists
}

func (s *Store) GetAllBookmarks() map[string]*Bookmark {
	return s.data.Bookmarks
}

func (s *Store) DeleteBookmark(alias string) error {
	delete(s.data.Bookmarks, alias)
	return s.save()
}

func (s *Store) UpdateUsage(alias string) error {
	if bookmark, exists := s.data.Bookmarks[alias]; exists {
		bookmark.UsedCount++
		bookmark.LastUsed = time.Now()
		return s.save()
	}
	return fmt.Errorf("bookmark not found: %s", alias)
}