package cmd

import (
	"fmt"
	"log"
	"strings"

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
	Use:   "labels [command]",
	Short: "Manage labels for Kubernetes contexts",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		lm, err := kube.NewLabelManager()
		if err != nil {
			log.Fatalf("Error creating label manager: %v", err)
		}

		defer func() {
			if err := lm.SaveLabels(); err != nil {
				log.Fatalf("Error saving labels: %v", err)
			}
		}()

		setLabelVal, _ := cmd.Flags().GetString("set")
		getLabelVal, _ := cmd.Flags().GetString("get")
		deleteLabelVal, _ := cmd.Flags().GetString("delete")
		contextNameVal, _ := cmd.Flags().GetString("context")

		if setLabelVal != "" {
			if contextNameVal == "" {
				log.Fatal("Context name is required to set a label.")
			}
			key, value := parseLabel(setLabelVal)
			lm.SetLabel(contextNameVal, key, value)
			fmt.Printf("Label '%s=%s' set for context '%s'.\n", key, value, contextNameVal)
			return
		} else if getLabelVal != "" {
			if contextNameVal == "" {
				log.Fatal("Context name is required to get a label.")
			}
			labels := lm.GetLabels(contextNameVal)
			if value, ok := labels[getLabelVal]; ok {
				fmt.Printf("Label '%s' for context '%s': %s\n", getLabelVal, contextNameVal, value)
			} else {
				fmt.Printf("Label '%s' not found for context '%s'.\n", getLabelVal, contextNameVal)
			}
			return
		} else if deleteLabelVal != "" {
			if contextNameVal == "" {
				log.Fatal("Context name is required to delete a label.")
			}
			lm.DeleteLabel(contextNameVal, deleteLabelVal)
			fmt.Printf("Label '%s' deleted from context '%s'.\n", deleteLabelVal, contextNameVal)
			return
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