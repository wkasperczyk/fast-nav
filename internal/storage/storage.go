package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// FuzzyMatch represents a fuzzy match result
type FuzzyMatch struct {
	Alias    string
	Bookmark *Bookmark
	Score    int
}

// FindFuzzyMatches finds bookmarks that match the given pattern
func (s *Store) FindFuzzyMatches(pattern string) []FuzzyMatch {
	var matches []FuzzyMatch
	pattern = strings.ToLower(pattern)
	
	for alias, bookmark := range s.data.Bookmarks {
		aliasLower := strings.ToLower(alias)
		score := calculateFuzzyScore(pattern, aliasLower)
		
		if score > 0 {
			matches = append(matches, FuzzyMatch{
				Alias:    alias,
				Bookmark: bookmark,
				Score:    score,
			})
		}
	}
	
	// Sort by score (descending) and then by usage count (descending)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Bookmark.UsedCount > matches[j].Bookmark.UsedCount
	})
	
	return matches
}

// calculateFuzzyScore calculates a score for how well the pattern matches the alias
func calculateFuzzyScore(pattern, alias string) int {
	if pattern == alias {
		return 1000 // Exact match gets highest score
	}
	
	if strings.Contains(alias, pattern) {
		// Substring match - higher score for matches at the beginning
		if strings.HasPrefix(alias, pattern) {
			return 800 + len(pattern)*10 // Prefix match
		}
		return 500 + len(pattern)*5 // Substring match
	}
	
	// Check for partial character matches
	score := 0
	patternIndex := 0
	
	for i, char := range alias {
		if patternIndex < len(pattern) && char == rune(pattern[patternIndex]) {
			score += 100
			if i == patternIndex {
				score += 50 // Bonus for sequential match
			}
			patternIndex++
		}
	}
	
	// Must match at least half the pattern characters
	if patternIndex < len(pattern)/2 {
		return 0
	}
	
	return score
}

// GetRecentlyUsed returns bookmarks sorted by last usage (most recent first)
func (s *Store) GetRecentlyUsed(limit int) []FuzzyMatch {
	var matches []FuzzyMatch
	
	for alias, bookmark := range s.data.Bookmarks {
		matches = append(matches, FuzzyMatch{
			Alias:    alias,
			Bookmark: bookmark,
			Score:    bookmark.UsedCount, // Use usage count as score for sorting
		})
	}
	
	// Sort by last used time (descending) and then by usage count (descending)
	sort.Slice(matches, func(i, j int) bool {
		if !matches[i].Bookmark.LastUsed.Equal(matches[j].Bookmark.LastUsed) {
			return matches[i].Bookmark.LastUsed.After(matches[j].Bookmark.LastUsed)
		}
		return matches[i].Bookmark.UsedCount > matches[j].Bookmark.UsedCount
	})
	
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	
	return matches
}

// GetSuggestions returns alias suggestions for typos/similar names
func (s *Store) GetSuggestions(input string, maxDistance int) []FuzzyMatch {
	var suggestions []FuzzyMatch
	inputLower := strings.ToLower(input)
	
	for alias, bookmark := range s.data.Bookmarks {
		aliasLower := strings.ToLower(alias)
		distance := levenshteinDistance(inputLower, aliasLower)
		
		// Only suggest if edit distance is reasonable
		if distance <= maxDistance && distance > 0 {
			score := 1000 - distance*100 // Higher score for lower distance
			if strings.HasPrefix(aliasLower, inputLower) {
				score += 200 // Bonus for prefix match
			}
			
			suggestions = append(suggestions, FuzzyMatch{
				Alias:    alias,
				Bookmark: bookmark,
				Score:    score,
			})
		}
	}
	
	// Sort by score (higher is better) and usage count
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Score != suggestions[j].Score {
			return suggestions[i].Score > suggestions[j].Score
		}
		return suggestions[i].Bookmark.UsedCount > suggestions[j].Bookmark.UsedCount
	})
	
	return suggestions
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	
	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}
	
	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}
	
	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}