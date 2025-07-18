package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch <context-name>",
	Short: "Switch to a different Kubernetes context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Switching to context %s (functionality not yet implemented).", args[0])
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
