package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all contexts...")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
