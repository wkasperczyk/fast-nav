package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(fn completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ fn completion bash > /etc/bash_completion.d/fn
  # macOS:
  $ fn completion bash > /usr/local/etc/bash_completion.d/fn

Zsh:
  # If shell completion is not already enabled, run once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  $ source <(fn completion zsh)

  # To load completions for each session, execute once:
  $ fn completion zsh > "${fpath[1]}/_fn"

Fish:
  $ fn completion fish | source

  # To load completions for each session, execute once:
  $ fn completion fish > ~/.config/fish/completions/fn.fish

PowerShell:
  PS> fn completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> fn completion powershell > fn.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}