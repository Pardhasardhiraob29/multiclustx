package cmd

import (
	"fmt"
	"log"

	"multiclustx/internal/kube"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kube.LoadKubeconfig("")
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		contexts := kube.GetContexts(config)

		for _, context := range contexts {
            fmt.Printf("Name: %s, Cluster: %s, User: %s, Namespace: %s\n", context.Name, context.Cluster, context.AuthInfo, context.Namespace)
        }
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}