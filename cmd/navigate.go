package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rethil/fn/internal/storage"
	"github.com/spf13/cobra"
)

var navigateCmd = &cobra.Command{
	Use:               "navigate <alias>",
	Aliases:           []string{"<alias>"},
	Short:             "Output path for navigation (used by shell function)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: aliasCompletionFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		// Try exact match first
		bookmark, exists := store.GetBookmark(alias)
		if exists {
			// Check if directory still exists
			if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
				return fmt.Errorf("directory no longer exists: %s", bookmark.Path)
			}

			// Update usage stats
			store.UpdateUsage(alias)

			// Output the path for shell to use
			fmt.Print(bookmark.Path)
			return nil
		}

		// Try fuzzy matching
		matches := store.FindFuzzyMatches(alias)
		if len(matches) == 0 {
			// Try smart suggestions for typos
			suggestions := store.GetSuggestions(alias, 3) // Allow up to 3 character edits
			if len(suggestions) > 0 {
				yellow := color.New(color.FgYellow)
				cyan := color.New(color.FgCyan)

				fmt.Fprintf(os.Stderr, "No exact match found for '%s'. Did you mean:\n\n", alias)
				for i, suggestion := range suggestions {
					if i >= 5 { // Limit to top 5 suggestions
						break
					}
					yellow.Fprintf(os.Stderr, "  %s", suggestion.Alias)
					fmt.Fprintf(os.Stderr, " -> ")
					cyan.Fprintf(os.Stderr, "%s\n", suggestion.Bookmark.Path)
				}
				fmt.Fprintf(os.Stderr, "\n")
			}
			return fmt.Errorf("no bookmarks found matching '%s'", alias)
		}

		// If we have exactly one match, use it
		if len(matches) == 1 {
			match := matches[0]

			// Check if directory still exists
			if _, err := os.Stat(match.Bookmark.Path); os.IsNotExist(err) {
				return fmt.Errorf("directory no longer exists: %s", match.Bookmark.Path)
			}

			// Update usage stats
			store.UpdateUsage(match.Alias)

			// Output the path for shell to use
			fmt.Print(match.Bookmark.Path)
			return nil
		}

		// Multiple matches - show them to the user
		yellow := color.New(color.FgYellow)
		cyan := color.New(color.FgCyan)

		fmt.Fprintf(os.Stderr, "Multiple matches found for '%s':\n\n", alias)
		for i, match := range matches {
			if i >= 5 { // Limit to top 5 matches
				break
			}
			yellow.Fprintf(os.Stderr, "  %s", match.Alias)
			fmt.Fprintf(os.Stderr, " -> ")
			cyan.Fprintf(os.Stderr, "%s\n", match.Bookmark.Path)
		}
		fmt.Fprintf(os.Stderr, "\nPlease use a more specific alias.\n")

		return fmt.Errorf("ambiguous match")
	},
}
