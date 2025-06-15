package cmd

import (
	"github.com/spf13/cobra"
	"github.com/rethil/fast-nav/internal/storage"
)

// aliasCompletionFunc provides completion for alias names
func aliasCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	store, err := storage.NewStore()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	
	bookmarks := store.GetAllBookmarks()
	var aliases []string
	
	for alias := range bookmarks {
		aliases = append(aliases, alias)
	}
	
	return aliases, cobra.ShellCompDirectiveNoFileComp
}