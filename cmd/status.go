package cmd

import (
	"fmt"
	"log"
	"os"

	"multiclustx/internal/kube"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check reachability and health of Kubernetes clusters",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kube.LoadKubeconfig("")
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		contexts := kube.GetContexts(config)

		fmt.Printf("%-30s %-15s %-20s\n", "NAME", "REACHABLE", "SERVER VERSION")
		fmt.Printf("%-30s %-15s %-20s\n", "----", "---------", "--------------")

		for _, context := range contexts {
			reachable := "No"
			serverVersion := "N/A"

			err := kube.PingTest(os.Getenv("KUBECONFIG"), context.Name)
			if err == nil {
				reachable = "Yes"
				version, err := kube.GetServerVersion(os.Getenv("KUBECONFIG"), context.Name)
				if err == nil {
					serverVersion = version
				}
			}
			fmt.Printf("%-30s %-15s %-20s\n", context.Name, reachable, serverVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
