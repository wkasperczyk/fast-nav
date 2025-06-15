package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rethil/fast-nav/internal/storage"
)

var pathCmd = &cobra.Command{
	Use:               "path <alias>",
	Short:             "Print path without navigating",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: aliasCompletionFunc,
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
		
		fmt.Println(bookmark.Path)
		return nil
	},
}