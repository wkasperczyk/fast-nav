package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fn",
	Short: "Fast navigation tool for bookmarking directories",
	Long: `fn is a command-line tool that allows you to bookmark directories 
and quickly navigate to them using short aliases.

Usage:
  fn save <alias>     Save current directory with an alias
  fn <alias>          Navigate to saved directory  
  fn list             List all saved aliases
  fn delete <alias>   Remove a saved alias
  fn path <alias>     Print path without navigating
  fn edit <alias>     Update existing alias to current directory
  fn cleanup          Remove bookmarks pointing to non-existent directories
  fn search <pattern> Find bookmarks by alias or path pattern`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(navigateCmd)
	rootCmd.AddCommand(pathCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(completionCmd)
}