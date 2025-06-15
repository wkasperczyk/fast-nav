package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"github.com/rethil/fn/internal/storage"
)

var recentCmd = &cobra.Command{
	Use:     "recent [index]",
	Aliases: []string{"r"},
	Short:   "Navigate to recently used bookmarks",
	Long: `Navigate to recently used bookmarks. 
If no index is provided, shows a list of recent bookmarks.
If an index is provided (1-9), navigates to that bookmark directly.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}
		
		recentBookmarks := store.GetRecentlyUsed(9) // Limit to 9 for single-digit selection
		
		if len(recentBookmarks) == 0 {
			fmt.Println("No bookmarks found")
			return nil
		}
		
		// If index provided, navigate directly
		if len(args) == 1 {
			index, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[0])
			}
			
			if index < 1 || index > len(recentBookmarks) {
				return fmt.Errorf("index out of range: %d (available: 1-%d)", index, len(recentBookmarks))
			}
			
			bookmark := recentBookmarks[index-1]
			
			// Check if directory still exists
			if _, err := os.Stat(bookmark.Bookmark.Path); os.IsNotExist(err) {
				return fmt.Errorf("directory no longer exists: %s", bookmark.Bookmark.Path)
			}
			
			// Update usage stats
			store.UpdateUsage(bookmark.Alias)
			
			// Output the path for shell to use
			fmt.Print(bookmark.Bookmark.Path)
			return nil
		}
		
		// Show list of recent bookmarks
		green := color.New(color.FgGreen)
		yellow := color.New(color.FgYellow)
		cyan := color.New(color.FgCyan)
		gray := color.New(color.FgHiBlack)
		
		fmt.Println("Recently used bookmarks:")
		fmt.Println()
		
		for i, bookmark := range recentBookmarks {
			index := i + 1
			
			// Format last used time
			var timeStr string
			now := time.Now()
			lastUsed := bookmark.Bookmark.LastUsed
			
			if lastUsed.IsZero() {
				timeStr = "never"
			} else {
				diff := now.Sub(lastUsed)
				if diff < time.Hour {
					timeStr = "< 1h ago"
				} else if diff < 24*time.Hour {
					timeStr = fmt.Sprintf("%dh ago", int(diff.Hours()))
				} else {
					timeStr = fmt.Sprintf("%dd ago", int(diff.Hours()/24))
				}
			}
			
			green.Printf("  %d. ", index)
			yellow.Printf("%-15s", bookmark.Alias)
			fmt.Printf(" → ")
			cyan.Printf("%-40s", bookmark.Bookmark.Path)
			gray.Printf(" (used %d times, %s)\n", bookmark.Bookmark.UsedCount, timeStr)
		}
		
		fmt.Println()
		fmt.Printf("Use 'fn recent <index>' to navigate directly (e.g., 'fn recent 1')\n")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}