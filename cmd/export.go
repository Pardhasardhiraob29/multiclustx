package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export results to a file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Export command not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
