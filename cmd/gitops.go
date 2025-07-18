package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gitopsCmd = &cobra.Command{
	Use:   "gitops",
	Short: "Show desired vs actual drift in GitOps deployments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GitOps sync functionality not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(gitopsCmd)
}
