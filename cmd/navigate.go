package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rethil/fn/internal/storage"
)

var navigateCmd = &cobra.Command{
	Use:   "navigate <alias>",
	Short: "Output path for navigation (used by shell function)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		bookmark, exists := store.GetBookmark(alias)
		if !exists {
			return fmt.Errorf("alias '%s' not found", alias)
		}
		
		// Check if directory still exists
		if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
			return fmt.Errorf("directory no longer exists: %s", bookmark.Path)
		}
		
		// Update usage stats
		store.UpdateUsage(alias)
		
		// Output the path for shell to use
		fmt.Print(bookmark.Path)
		return nil
	},
}