package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/rethil/fn/internal/storage"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <alias>",
	Short: "Remove a saved alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		_, exists := store.GetBookmark(alias)
		if !exists {
			return fmt.Errorf("alias '%s' not found", alias)
		}
		
		// Confirmation prompt
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete '%s'?", alias),
			Default: false,
		}
		
		err = survey.AskOne(prompt, &confirm)
		if err != nil {
			return err
		}
		
		if !confirm {
			fmt.Println("Cancelled.")
			return nil
		}
		
		err = store.DeleteBookmark(alias)
		if err != nil {
			return fmt.Errorf("failed to delete bookmark: %w", err)
		}
		
		fmt.Printf("✓ Deleted '%s'\n", alias)
		return nil
	},
}