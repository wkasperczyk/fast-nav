package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rethil/fn/internal/storage"
)

// Integration tests for CLI commands
// These tests focus on core functionality without mocking cobra's execution

func TestCLIIntegration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fn-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test directory to navigate to
	testDir := filepath.Join(tempDir, "test-nav-dir")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test nav dir: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	originalPwd, _ := os.Getwd()
	
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Chdir(originalPwd)
	}()

	os.Setenv("HOME", tempDir)

	// Change to test directory
	err = os.Chdir(testDir)
	if err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	// Test 1: Save command functionality
	t.Run("SaveBookmark", func(t *testing.T) {
		err := saveCmd.RunE(saveCmd, []string{"test-save"})
		if err != nil {
			t.Fatalf("Save command failed: %v", err)
		}

		// Verify bookmark was saved
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		bookmark, exists := store.GetBookmark("test-save")
		if !exists {
			t.Error("Bookmark was not saved")
		}

		if bookmark.Path != testDir {
			t.Errorf("Expected path %s, got %s", testDir, bookmark.Path)
		}
	})

	// Test 2: List command functionality
	t.Run("ListBookmarks", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		// Add test bookmarks
		err = store.SaveBookmark("list-test1", "/path/one")
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		err = store.SaveBookmark("list-test2", "/path/two")
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		// Test that GetAllBookmarks works
		bookmarks := store.GetAllBookmarks()
		if len(bookmarks) < 2 {
			t.Errorf("Expected at least 2 bookmarks, got %d", len(bookmarks))
		}

		if _, exists := bookmarks["list-test1"]; !exists {
			t.Error("Expected bookmark 'list-test1' to exist")
		}

		if _, exists := bookmarks["list-test2"]; !exists {
			t.Error("Expected bookmark 'list-test2' to exist")
		}
	})

	// Test 3: Delete command functionality
	t.Run("DeleteBookmark", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		// Add test bookmark
		err = store.SaveBookmark("delete-test", "/path/to/delete")
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		// Delete bookmark
		err = store.DeleteBookmark("delete-test")
		if err != nil {
			t.Fatalf("Failed to delete bookmark: %v", err)
		}

		// Verify bookmark was deleted
		_, exists := store.GetBookmark("delete-test")
		if exists {
			t.Error("Bookmark should have been deleted")
		}
	})

	// Test 4: Path command functionality
	t.Run("PathCommand", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		testPath := "/path/to/test"
		err = store.SaveBookmark("path-test", testPath)
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		// Test that GetBookmark returns correct path
		bookmark, exists := store.GetBookmark("path-test")
		if !exists {
			t.Error("Bookmark should exist")
		}

		if bookmark.Path != testPath {
			t.Errorf("Expected path %s, got %s", testPath, bookmark.Path)
		}
	})

	// Test 5: Edit command functionality
	t.Run("EditBookmark", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		// Add bookmark with original path
		err = store.SaveBookmark("edit-test", "/original/path")
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		// Update bookmark to current directory (testDir)
		err = store.SaveBookmark("edit-test", testDir)
		if err != nil {
			t.Fatalf("Failed to update bookmark: %v", err)
		}

		// Verify bookmark was updated
		bookmark, exists := store.GetBookmark("edit-test")
		if !exists {
			t.Error("Bookmark should still exist")
		}

		if bookmark.Path != testDir {
			t.Errorf("Expected updated path %s, got %s", testDir, bookmark.Path)
		}
	})

	// Test 6: Navigation usage tracking
	t.Run("NavigationUsageTracking", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		// Add bookmark
		err = store.SaveBookmark("usage-test", "/path/for/usage")
		if err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}

		// Update usage (simulate navigation)
		err = store.UpdateUsage("usage-test")
		if err != nil {
			t.Fatalf("Failed to update usage: %v", err)
		}

		// Verify usage was updated
		bookmark, exists := store.GetBookmark("usage-test")
		if !exists {
			t.Error("Bookmark should exist")
		}

		if bookmark.UsedCount != 1 {
			t.Errorf("Expected UsedCount 1, got %d", bookmark.UsedCount)
		}
	})

	// Test 7: Cleanup functionality 
	t.Run("CleanupInvalidBookmarks", func(t *testing.T) {
		store, err := storage.NewStore()
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}

		// Add valid bookmark (existing directory)
		err = store.SaveBookmark("cleanup-valid", testDir)
		if err != nil {
			t.Fatalf("Failed to save valid bookmark: %v", err)
		}

		// Add invalid bookmark (non-existent directory)
		err = store.SaveBookmark("cleanup-invalid", "/definitely/does/not/exist")
		if err != nil {
			t.Fatalf("Failed to save invalid bookmark: %v", err)
		}

		// Get initial bookmarks
		initialBookmarks := store.GetAllBookmarks()

		// Test cleanup logic by manually checking what would be cleaned
		invalidCount := 0
		for _, bookmark := range initialBookmarks {
			if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
				invalidCount++
			}
		}

		if invalidCount == 0 {
			t.Error("Expected to find at least one invalid bookmark for cleanup test")
		}

		// Verify valid bookmark still exists
		_, exists := store.GetBookmark("cleanup-valid")
		if !exists {
			t.Error("Valid bookmark should still exist")
		}
	})
}

func TestAliasValidation(t *testing.T) {
	tests := []struct {
		alias     string
		expected  bool
		desc      string
	}{
		{"valid-alias", true, "valid alias with dash"},
		{"valid_alias", true, "valid alias with underscore"},
		{"ValidAlias123", true, "valid alias with mixed case and numbers"},
		{"save", false, "reserved word"},
		{"list", false, "reserved word"},
		{"", false, "empty alias"},
		{"alias with spaces", false, "alias with spaces"},
		{"alias@invalid", false, "alias with special characters"},
		{strings.Repeat("a", 51), false, "alias too long"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := isValidAlias(tt.alias)
			if result != tt.expected {
				t.Errorf("isValidAlias(%q) = %v, expected %v", tt.alias, result, tt.expected)
			}
		})
	}
}

func TestStorageIntegration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fn-storage-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test storage creation and basic operations
	store, err := storage.NewStore()
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Test saving and retrieving bookmarks
	testCases := []struct {
		alias string
		path  string
	}{
		{"home", "/home/user"},
		{"work", "/work/project"},
		{"docs", "/home/user/documents"},
	}

	for _, tc := range testCases {
		err = store.SaveBookmark(tc.alias, tc.path)
		if err != nil {
			t.Fatalf("Failed to save bookmark %s: %v", tc.alias, err)
		}

		bookmark, exists := store.GetBookmark(tc.alias)
		if !exists {
			t.Errorf("Bookmark %s should exist", tc.alias)
		}

		if bookmark.Path != tc.path {
			t.Errorf("Expected path %s for %s, got %s", tc.path, tc.alias, bookmark.Path)
		}
	}

	// Test getting all bookmarks
	allBookmarks := store.GetAllBookmarks()
	if len(allBookmarks) != len(testCases) {
		t.Errorf("Expected %d bookmarks, got %d", len(testCases), len(allBookmarks))
	}

	// Test usage tracking
	err = store.UpdateUsage("home")
	if err != nil {
		t.Fatalf("Failed to update usage: %v", err)
	}

	bookmark, exists := store.GetBookmark("home")
	if !exists {
		t.Error("Bookmark should exist")
	}

	if bookmark.UsedCount != 1 {
		t.Errorf("Expected UsedCount 1, got %d", bookmark.UsedCount)
	}

	// Test deletion
	err = store.DeleteBookmark("work")
	if err != nil {
		t.Fatalf("Failed to delete bookmark: %v", err)
	}

	_, exists = store.GetBookmark("work")
	if exists {
		t.Error("Bookmark should have been deleted")
	}
}