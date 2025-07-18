package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels for Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Labels command not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(labelsCmd)
}
