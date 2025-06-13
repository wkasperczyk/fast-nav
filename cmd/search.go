package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/rethil/fn/internal/storage"
)

var searchCmd = &cobra.Command{
	Use:   "search <pattern>",
	Short: "Find bookmarks by alias or path pattern",
	Long:  `Search for bookmarks matching a pattern in either alias names or directory paths.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := strings.ToLower(args[0])
		
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		bookmarks := store.GetAllBookmarks()
		var matches []string
		
		for alias, bookmark := range bookmarks {
			aliasLower := strings.ToLower(alias)
			pathLower := strings.ToLower(bookmark.Path)
			
			if strings.Contains(aliasLower, pattern) || strings.Contains(pathLower, pattern) {
				matches = append(matches, alias)
			}
		}
		
		if len(matches) == 0 {
			color.Red("No bookmarks found matching '%s'", pattern)
			return nil
		}
		
		color.Cyan("🔍 Found %d bookmark(s) matching '%s':", len(matches), pattern)
		for _, alias := range matches {
			bookmark := bookmarks[alias]
			fmt.Printf("📍 %-12s → %s (used %d times)\n", alias, bookmark.Path, bookmark.UsedCount)
		}
		
		return nil
	},
}