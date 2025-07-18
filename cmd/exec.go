package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"multiclustx/internal/executor"
	"multiclustx/internal/kube"
	"multiclustx/pkg/types"

	"github.com/spf13/cobra"
)

var (
	allClusters bool
	labelFilter string
)

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Execute kubectl-like commands across multiple clusters",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kube.LoadKubeconfig("")
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		contexts := kube.GetContexts(config)

		lm, err := kube.NewLabelManager()
		if err != nil {
			log.Fatalf("Error creating label manager: %v", err)
		}

		var contextsToExecute []types.ContextInfo
		if allClusters {
			contextsToExecute = contexts
		} else if labelFilter != "" {
			contextsToExecute = kube.FilterContextsByLabel(contexts, lm, labelFilter)
		} else {
			// If no flags are specified, execute on the current context (not yet implemented)
			log.Fatal("Please specify --all-clusters or --label.")
		}

		if len(contextsToExecute) == 0 {
			fmt.Println("No contexts found to execute on.")
			return
		}

		for _, context := range contextsToExecute {
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatalf("Error getting user home directory: %v", err)
				}
				kubeconfigPath = filepath.Join(home, ".kube", "config")
			}

			fmt.Printf("\n--- Executing on context: %s ---\n", context.Name)
			stdout, stderr, err := executor.ExecuteKubectlCommand(kubeconfigPath, context.Name, args)
			if err != nil {
				fmt.Printf("Error: %v\nStderr: %s\n", err, stderr)
			} else {
				fmt.Printf("Stdout:\n%s\n", stdout)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().BoolVarP(&allClusters, "all-clusters", "a", false, "Run command across all clusters")
	execCmd.Flags().StringVarP(&labelFilter, "label", "l", "", "Run command on clusters with the specified label")
}
