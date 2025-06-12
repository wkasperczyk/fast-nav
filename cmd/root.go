package cmd

import (
	"fmt"
	"os"

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
  fn edit <alias>     Update existing alias to current directory`,
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
}