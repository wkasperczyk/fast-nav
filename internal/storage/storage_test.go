package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fn-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	store, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore() failed: %v", err)
	}

	if store == nil {
		t.Fatal("NewStore() returned nil store")
	}

	// Check that config directory was created
	configDir := filepath.Join(tempDir, ".fn")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", configDir)
	}

	// Check that bookmarks file was created
	bookmarksFile := filepath.Join(configDir, "bookmarks.json")
	if _, err := os.Stat(bookmarksFile); os.IsNotExist(err) {
		t.Errorf("Bookmarks file was not created: %s", bookmarksFile)
	}
}

func TestSaveBookmark(t *testing.T) {
	store := setupTestStore(t)
	
	err := store.SaveBookmark("test", "/tmp/test")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	bookmark, exists := store.GetBookmark("test")
	if !exists {
		t.Fatal("Bookmark was not saved")
	}

	if bookmark.Path != "/tmp/test" {
		t.Errorf("Expected path '/tmp/test', got '%s'", bookmark.Path)
	}

	if bookmark.UsedCount != 0 {
		t.Errorf("Expected UsedCount 0, got %d", bookmark.UsedCount)
	}
}

func TestUpdateExistingBookmark(t *testing.T) {
	store := setupTestStore(t)
	
	// Save initial bookmark
	err := store.SaveBookmark("test", "/tmp/test1")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	// Update existing bookmark
	err = store.SaveBookmark("test", "/tmp/test2")
	if err != nil {
		t.Fatalf("SaveBookmark() update failed: %v", err)
	}

	bookmark, exists := store.GetBookmark("test")
	if !exists {
		t.Fatal("Bookmark does not exist")
	}

	if bookmark.Path != "/tmp/test2" {
		t.Errorf("Expected updated path '/tmp/test2', got '%s'", bookmark.Path)
	}
}

func TestGetBookmark(t *testing.T) {
	store := setupTestStore(t)
	
	// Test non-existent bookmark
	_, exists := store.GetBookmark("nonexistent")
	if exists {
		t.Error("Expected bookmark to not exist")
	}

	// Save and retrieve bookmark
	err := store.SaveBookmark("test", "/tmp/test")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	bookmark, exists := store.GetBookmark("test")
	if !exists {
		t.Fatal("Expected bookmark to exist")
	}

	if bookmark.Path != "/tmp/test" {
		t.Errorf("Expected path '/tmp/test', got '%s'", bookmark.Path)
	}
}

func TestDeleteBookmark(t *testing.T) {
	store := setupTestStore(t)
	
	// Save bookmark
	err := store.SaveBookmark("test", "/tmp/test")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	// Delete bookmark
	err = store.DeleteBookmark("test")
	if err != nil {
		t.Fatalf("DeleteBookmark() failed: %v", err)
	}

	// Verify it's deleted
	_, exists := store.GetBookmark("test")
	if exists {
		t.Error("Expected bookmark to be deleted")
	}
}

func TestUpdateUsage(t *testing.T) {
	store := setupTestStore(t)
	
	// Save bookmark
	err := store.SaveBookmark("test", "/tmp/test")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	// Update usage
	err = store.UpdateUsage("test")
	if err != nil {
		t.Fatalf("UpdateUsage() failed: %v", err)
	}

	bookmark, exists := store.GetBookmark("test")
	if !exists {
		t.Fatal("Bookmark does not exist")
	}

	if bookmark.UsedCount != 1 {
		t.Errorf("Expected UsedCount 1, got %d", bookmark.UsedCount)
	}

	// Test non-existent bookmark
	err = store.UpdateUsage("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent bookmark")
	}
}

func TestGetAllBookmarks(t *testing.T) {
	store := setupTestStore(t)
	
	// Save multiple bookmarks
	err := store.SaveBookmark("test1", "/tmp/test1")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	err = store.SaveBookmark("test2", "/tmp/test2")
	if err != nil {
		t.Fatalf("SaveBookmark() failed: %v", err)
	}

	bookmarks := store.GetAllBookmarks()
	
	// Check we have the expected bookmarks
	if _, exists := bookmarks["test1"]; !exists {
		t.Error("Expected bookmark 'test1' to exist")
	}
	
	if _, exists := bookmarks["test2"]; !exists {
		t.Error("Expected bookmark 'test2' to exist")
	}

	if bookmarks["test1"].Path != "/tmp/test1" {
		t.Errorf("Expected test1 path '/tmp/test1', got '%s'", bookmarks["test1"].Path)
	}

	if bookmarks["test2"].Path != "/tmp/test2" {
		t.Errorf("Expected test2 path '/tmp/test2', got '%s'", bookmarks["test2"].Path)
	}
}

// setupTestStore creates a temporary store for testing
func setupTestStore(t *testing.T) *Store {
	tempDir, err := os.MkdirTemp("", "fn-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Clean up after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	store, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore() failed: %v", err)
	}

	return store
}