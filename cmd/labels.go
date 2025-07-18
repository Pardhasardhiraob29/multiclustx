package cmd

import (
	"fmt"
	"log"

	"multiclustx/internal/kube"

	"github.com/spf13/cobra"
)

var ( 
	setLabel string
	getLabel string
	deleteLabel string
	contextName string
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels for Kubernetes contexts",
	Run: func(cmd *cobra.Command, args []string) {
		lm, err := kube.NewLabelManager()
		if err != nil {
			log.Fatalf("Error creating label manager: %v", err)
		}

		if setLabel != "" {
			if contextName == "" {
				log.Fatal("Context name is required to set a label.")
			}
			key, value := parseLabel(setLabel)
			lm.SetLabel(contextName, key, value)
			fmt.Printf("Label '%s=%s' set for context '%s'.\n", key, value, contextName)
		} else if getLabel != "" {
			if contextName == "" {
				log.Fatal("Context name is required to get a label.")
			}
			labels := lm.GetLabels(contextName)
			if value, ok := labels[getLabel]; ok {
				fmt.Printf("Label '%s' for context '%s': %s\n", getLabel, contextName, value)
			} else {
				fmt.Printf("Label '%s' not found for context '%s'.\n", getLabel, contextName)
			}
		} else if deleteLabel != "" {
			if contextName == "" {
				log.Fatal("Context name is required to delete a label.")
			}
			lm.DeleteLabel(contextName, deleteLabel)
			fmt.Printf("Label '%s' deleted from context '%s'.\n", deleteLabel, contextName)
		} else {
			// List all labels for all contexts
			allLabels := lm.GetAllContextLabels()
			if len(allLabels) == 0 {
				fmt.Println("No labels found.")
				return
			}
			fmt.Println("All Labels:")
			for ctx, labels := range allLabels {
				fmt.Printf("  Context: %s\n", ctx)
				for key, value := range labels {
					fmt.Printf("    %s: %s\n", key, value)
				}
			}
		}

		if err := lm.SaveLabels(); err != nil {
			log.Fatalf("Error saving labels: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(labelsCmd)

	labelsCmd.Flags().StringVarP(&setLabel, "set", "s", "", "Set a label for a context (e.g., 'env=prod')")
	labelsCmd.Flags().StringVarP(&getLabel, "get", "g", "", "Get a label for a context by key")
	labelsCmd.Flags().StringVarP(&deleteLabel, "delete", "d", "", "Delete a label from a context by key")
	labelsCmd.Flags().StringVarP(&contextName, "context", "c", "", "Specify the context name")
}

func parseLabel(label string) (string, string) {
	parts := strings.SplitN(label, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return label, ""
}