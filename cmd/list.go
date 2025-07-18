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

				fmt.Printf("%-30s %-30s %-30s %-30s\n", "NAME", "CLUSTER", "USER", "NAMESPACE")
		fm t.Printf("%-30s %-30s %-30s %-30s\n", "----", "-------", "----", "---------")
		for _, context := range contexts {
			fm t.Printf("%-30s %-30s %-30s %-30s\n", context.Name, context.Cluster, context.AuthInfo, context.Namespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}