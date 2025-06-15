package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rethil/fast-nav/internal/storage"
)

var editCmd = &cobra.Command{
	Use:               "edit <alias>",
	Short:             "Update existing alias to current directory", 
	Long:              `Update an existing bookmark alias to point to the current directory.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: aliasCompletionFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		// Check if alias exists
		_, exists := store.GetBookmark(alias)
		if !exists {
			return fmt.Errorf("alias '%s' does not exist", alias)
		}
		
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		
		// Update the existing bookmark
		err = store.SaveBookmark(alias, currentDir)
		if err != nil {
			return fmt.Errorf("failed to update bookmark: %w", err)
		}
		
		fmt.Printf("Updated alias '%s' to: %s\n", alias, currentDir)
		return nil
	},
}