package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra" 
	"github.com/rethil/fn/internal/storage"
)

var completionCmd = &cobra.Command{
	Use:    "completion <word>",
	Short:  "Generate tab completion for aliases",
	Long:   `Generate tab completion suggestions for bookmark aliases.`,
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentWord := args[0]
		
		store, err := storage.NewStore()
		if err != nil {
			return err
		}
		
		bookmarks := store.GetAllBookmarks()
		
		for alias := range bookmarks {
			if strings.HasPrefix(alias, currentWord) {
				fmt.Println(alias)
			}
		}
		
		return nil
	},
}