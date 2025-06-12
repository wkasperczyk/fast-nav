package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/rethil/fn/internal/storage"
)

var saveCmd = &cobra.Command{
	Use:   "save <alias>",
	Short: "Save current directory with an alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		
		// Validate alias
		if !isValidAlias(alias) {
			return fmt.Errorf("invalid alias: use only alphanumeric characters, dash, and underscore")
		}
		
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		
		// Save bookmark
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		err = store.SaveBookmark(alias, currentDir)
		if err != nil {
			return fmt.Errorf("failed to save bookmark: %w", err)
		}
		
		fmt.Printf("✓ Saved '%s' → %s\n", alias, currentDir)
		return nil
	},
}

func isValidAlias(alias string) bool {
	// Check reserved words
	reserved := []string{"save", "list", "delete", "edit", "path", "help", "navigate"}
	for _, word := range reserved {
		if alias == word {
			return false
		}
	}
	
	// Check format: alphanumeric + dash/underscore, max 50 chars
	if len(alias) > 50 {
		return false
	}
	
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, alias)
	return matched
}