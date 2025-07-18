package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multiclustx",
	Short: "A CLI tool to manage multiple Kubernetes clusters.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to multiclustx!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
