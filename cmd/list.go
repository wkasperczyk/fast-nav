package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/rethil/fn/internal/storage"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		bookmarks := store.GetAllBookmarks()
		if len(bookmarks) == 0 {
			fmt.Println("No bookmarks saved yet. Use 'fn save <alias>' to create one.")
			return nil
		}
		
		for alias, bookmark := range bookmarks {
			// Check if directory still exists
			exists := true
			if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
				exists = false
			}
			
			if exists {
				color.Green("📍 %-12s → %s (used %d times)", alias, bookmark.Path, bookmark.UsedCount)
			} else {
				color.Red("❌ %-12s → %s (MISSING - used %d times)", alias, bookmark.Path, bookmark.UsedCount)
			}
		}
		
		return nil
	},
}