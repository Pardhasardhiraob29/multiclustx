package cmd

import (
	"fmt"
	"log"
	"os"

	"multiclustx/internal/kube"

	"github.com/olekukonko/tablewriter"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Cluster", "User", "Namespace"})

		for _, context := range contexts {
			row := []string{context.Name, context.Cluster, context.AuthInfo, context.Namespace}
			table.Append(row)
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}