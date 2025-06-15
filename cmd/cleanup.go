package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/rethil/fast-nav/internal/storage"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove bookmarks pointing to non-existent directories",
	Long:  `Remove all bookmarks that point to directories that no longer exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		bookmarks := store.GetAllBookmarks()
		var removed []string
		
		for alias, bookmark := range bookmarks {
			if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
				err := store.DeleteBookmark(alias)
				if err != nil {
					return fmt.Errorf("failed to delete bookmark '%s': %w", alias, err)
				}
				removed = append(removed, alias)
			}
		}
		
		if len(removed) == 0 {
			color.Green("✓ All bookmarks are valid - no cleanup needed")
		} else {
			color.Yellow("🧹 Cleaned up %d invalid bookmarks:", len(removed))
			for _, alias := range removed {
				fmt.Printf("  • %s\n", alias)
			}
		}
		
		return nil
	},
}