package cmd

import (
	"fmt"
	"log"
	"os"

	"multiclustx/internal/kube"
	"multiclustx/internal/rbac"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit Kubernetes contexts for RBAC capabilities",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kube.LoadKubeconfig("")
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		contexts := kube.GetContexts(config)

		fmt.Printf("%-30s %-15s %-20s\n", "CONTEXT", "RBAC STATUS", "RULES")
		fmt.Printf("%-30s %-15s %-20s\n", "-------", "-----------", "-----")

		for _, context := range contexts {
			rulesReview, err := rbac.CheckRBAC(os.Getenv("KUBECONFIG"), context.Name)
			if err != nil {
				fmt.Printf("%-30s %-15s %-20s\n", context.Name, "Error", err.Error())
				continue
			}

			status := "OK"
			rules := ""
			for _, rule := range rulesReview.Status.Rules {
				rules += fmt.Sprintf("Verbs: %v, Resources: %v\n", rule.Verbs, rule.Resources)
			}
			fmt.Printf("%-30s %-15s %-20s\n", context.Name, status, rules)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
