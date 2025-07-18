package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"multiclustx/internal/executor"
	"multiclustx/internal/kube"
	"multiclustx/pkg/types"

	"github.com/spf13/cobra"
	sigsyaml "sigs.k8s.io/yaml"
)

var (
	allClusters bool
	labelFilter string
	outputFormat string
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

		type Result struct {
			Context string `json:"context" yaml:"context"`
			Stdout  string `json:"stdout" yaml:"stdout"`
			Stderr  string `json:"stderr" yaml:"stderr"`
			Error   string `json:"error" yaml:"error"`
		}

		var results []Result

		for _, context := range contextsToExecute {
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatalf("Error getting user home directory: %v", err)
				}
				kubeconfigPath = filepath.Join(home, ".kube", "config")
			}

			stdout, stderr, err := executor.ExecuteKubectlCommand(kubeconfigPath, context.Name, args)
			errorStr := ""
			if err != nil {
				errorStr = err.Error()
			}
			results = append(results, Result{
				Context: context.Name,
				Stdout:  stdout,
				Stderr:  stderr,
				Error:   errorStr,
			})
		}

		switch outputFormat {
		case "json":
			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling JSON: %v", err)
			}
			fmt.Println(string(jsonData))
		case "yaml":
			yamlData, err := sigsyaml.Marshal(results)
			if err != nil {
				log.Fatalf("Error marshaling YAML: %v", err)
			}
			fmt.Println(string(yamlData))
		case "table", "":
			fmt.Printf("%-30s %-50s %-50s %-50s\n", "CONTEXT", "STDOUT", "STDERR", "ERROR")
			fmt.Printf("%-30s %-50s %-50s %-50s\n", "-------", "------", "------", "-----")
			for _, result := range results {
				fmt.Printf("%-30s %-50s %-50s %-50s\n", result.Context, truncate(result.Stdout, 47), truncate(result.Stderr, 47), truncate(result.Error, 47))
			}
		default:
			log.Fatalf("Unsupported output format: %s", outputFormat)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().BoolVarP(&allClusters, "all-clusters", "a", false, "Run command across all clusters")
	execCmd.Flags().StringVarP(&labelFilter, "label", "l", "", "Run command on clusters with the specified label")
	execCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, yaml)")
}

func truncate(s string, i int) string {
	if len(s) <= i {
		return s
	}
	return s[0:i] + "..."
}