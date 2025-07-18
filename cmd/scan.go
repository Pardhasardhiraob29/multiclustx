package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"multiclustx/internal/kube"
	"multiclustx/internal/scanner"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Kubernetes secrets for sensitive information",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kube.LoadKubeconfig("")
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		contexts := kube.GetContexts(config)

		fmt.Printf("%-30s %-15s\n", "CONTEXT", "FOUND TOKENS")
		fmt.Printf("%-30s %-15s\n", "-------", "------------")

		for _, context := range contexts {
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatalf("Error getting user home directory: %v", err)
				}
				kubeconfigPath = filepath.Join(home, ".kube", "config")
			}

			foundTokens, err := scanner.ScanSecretsForTokens(kubeconfigPath, context.Name)
			if err != nil {
				fmt.Printf("%-30s %-15s\n", context.Name, fmt.Sprintf("Error: %v", err))
				continue
			}

			if len(foundTokens) > 0 {
				for _, token := range foundTokens {
					fmt.Printf("%-30s %-15s\n", context.Name, token)
				}
			} else {
				fmt.Printf("%-30s %-15s\n", context.Name, "None")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
