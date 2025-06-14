package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func BenchmarkNewStore(b *testing.B) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fn-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Setenv("HOME", tempDir)
		store, err := NewStore()
		if err != nil {
			b.Fatalf("NewStore() failed: %v", err)
		}
		if store == nil {
			b.Fatal("NewStore() returned nil")
		}
	}
}

func BenchmarkSaveBookmark(b *testing.B) {
	store := setupBenchStore(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("bench-save-%d", i)
		err := store.SaveBookmark(alias, "/tmp/test/path")
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}
}

func BenchmarkSaveBookmarkUpdate(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-save the bookmark that will be updated
	err := store.SaveBookmark("update-test", "/original/path")
	if err != nil {
		b.Fatalf("Initial SaveBookmark() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newPath := fmt.Sprintf("/updated/path/%d", i)
		err := store.SaveBookmark("update-test", newPath)
		if err != nil {
			b.Fatalf("SaveBookmark() update failed: %v", err)
		}
	}
}

func BenchmarkGetBookmark(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with test data
	for i := 0; i < 100; i++ {
		alias := fmt.Sprintf("bench-get-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("bench-get-%d", i%100)
		_, exists := store.GetBookmark(alias)
		if !exists {
			b.Fatalf("Expected bookmark %s to exist", alias)
		}
	}
}

func BenchmarkGetBookmarkMiss(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with some test data
	for i := 0; i < 10; i++ {
		alias := fmt.Sprintf("bench-existing-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("nonexistent-%d", i)
		_, exists := store.GetBookmark(alias)
		if exists {
			b.Fatalf("Expected bookmark %s to not exist", alias)
		}
	}
}

func BenchmarkGetAllBookmarks(b *testing.B) {
	store := setupBenchStore(b)

	// Test with different numbers of bookmarks
	benchmarks := []struct {
		name  string
		count int
	}{
		{"10", 10},
		{"100", 100},
		{"1000", 1000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Pre-populate with test data
			for i := 0; i < bm.count; i++ {
				alias := fmt.Sprintf("bench-all-%d", i)
				err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
				if err != nil {
					b.Fatalf("SaveBookmark() failed: %v", err)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bookmarks := store.GetAllBookmarks()
				if len(bookmarks) < bm.count {
					b.Fatalf("Expected at least %d bookmarks, got %d", bm.count, len(bookmarks))
				}
			}
		})
	}
}

func BenchmarkDeleteBookmark(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with bookmarks to delete
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("bench-delete-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("bench-delete-%d", i)
		err := store.DeleteBookmark(alias)
		if err != nil {
			b.Fatalf("DeleteBookmark() failed: %v", err)
		}
	}
}

func BenchmarkUpdateUsage(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with test data
	for i := 0; i < 100; i++ {
		alias := fmt.Sprintf("bench-usage-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := fmt.Sprintf("bench-usage-%d", i%100)
		err := store.UpdateUsage(alias)
		if err != nil {
			b.Fatalf("UpdateUsage() failed: %v", err)
		}
	}
}

func BenchmarkLoadStore(b *testing.B) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fn-bench-load-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Create a store with test data
	store, err := NewStore()
	if err != nil {
		b.Fatalf("NewStore() failed: %v", err)
	}

	// Populate with test data
	for i := 0; i < 100; i++ {
		alias := fmt.Sprintf("load-test-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Force reload by creating a new store instance
		_, err := NewStore()
		if err != nil {
			b.Fatalf("NewStore() failed: %v", err)
		}
	}
}

func BenchmarkSaveStore(b *testing.B) {
	store := setupBenchStore(b)

	// Populate with test data
	now := time.Now()
	for i := 0; i < 100; i++ {
		alias := fmt.Sprintf("save-test-%d", i)
		store.data.Bookmarks[alias] = &Bookmark{
			Path:      fmt.Sprintf("/path/%d", i),
			Created:   now,
			UsedCount: i,
			LastUsed:  now,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := store.save()
		if err != nil {
			b.Fatalf("save() failed: %v", err)
		}
	}
}

func BenchmarkConcurrentRead(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with test data
	for i := 0; i < 100; i++ {
		alias := fmt.Sprintf("concurrent-read-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			alias := fmt.Sprintf("concurrent-read-%d", i%100)
			_, exists := store.GetBookmark(alias)
			if !exists {
				b.Fatalf("Expected bookmark %s to exist", alias)
			}
			i++
		}
	})
}

func BenchmarkConcurrentWrite(b *testing.B) {
	store := setupBenchStore(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			alias := fmt.Sprintf("concurrent-write-%d", i)
			err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
			if err != nil {
				b.Fatalf("SaveBookmark() failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkMixedOperations(b *testing.B) {
	store := setupBenchStore(b)

	// Pre-populate with some initial data
	for i := 0; i < 50; i++ {
		alias := fmt.Sprintf("mixed-%d", i)
		err := store.SaveBookmark(alias, fmt.Sprintf("/path/%d", i))
		if err != nil {
			b.Fatalf("SaveBookmark() failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0: // Save new bookmark
			alias := fmt.Sprintf("mixed-new-%d", i)
			store.SaveBookmark(alias, fmt.Sprintf("/new/path/%d", i))
		case 1: // Get existing bookmark
			alias := fmt.Sprintf("mixed-%d", i%50)
			store.GetBookmark(alias)
		case 2: // Update usage
			alias := fmt.Sprintf("mixed-%d", i%50)
			store.UpdateUsage(alias)
		case 3: // Get all bookmarks
			store.GetAllBookmarks()
		}
	}
}

// setupBenchStore creates a temporary store for benchmarking
func setupBenchStore(b *testing.B) *Store {
	tempDir, err := os.MkdirTemp("", "fn-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// Clean up after benchmark
	b.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	b.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	store, err := NewStore()
	if err != nil {
		b.Fatalf("NewStore() failed: %v", err)
	}

	return store
}