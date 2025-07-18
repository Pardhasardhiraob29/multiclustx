package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To generate the desired shell's completion script, use:

Bash:
  $ source <(multiclustx completion bash)

  # To load completions for each session, add this to your bashrc:
  $ echo "source <(multiclustx completion bash)" >> ~/.bashrc

Zsh:
  # If shell completion is not already enabled in your environment, enable it.
  # You will likely need to add the following to your .zshrc:
  #
  #   autoload -Uz compinit
  #   compinit
  #
  # Afterwards, simply run:
  $ multiclustx completion zsh > ~/.zsh/_multiclustx
  #
  # Then, add the following to your .zshrc:
  #
  #   fpath=(~/.zsh $fpath)
  #   compinit
  #
  # You will need to force a reload of your .zshrc:
  #
  #   source ~/.zshrc

Fish:
  $ multiclustx completion fish | source

  # To load completions for each session, add this to your fish.config:
  $ multiclustx completion fish > ~/.config/fish/completions/multiclustx.fish

Powershell:
  PS> multiclustx completion powershell | Out-String | Invoke-Expression

  # To load completions for each session, add this to your powershell profile:
  PS> multiclustx completion powershell | Out-String | Add-Content -Path $PROFILE
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowershellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
