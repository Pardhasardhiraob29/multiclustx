package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute kubectl-like commands across multiple clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exec command not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
