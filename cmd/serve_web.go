package cmd

import (
	"fmt"

	"multiclustx/internal/webserver"

	"github.com/spf13/cobra"
)

var serveWebCmd = &cobra.Command{
	Use:   "serve-web",
	Short: "Serve the MultiClustX web UI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting MultiClustX web UI...")
		webserver.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveWebCmd)
}
