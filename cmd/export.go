package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	sigsyaml "sigs.k8s.io/yaml"
)

var (
	exportOutputFormat string
	exportFilePath string
)

var exportCmd = &cobra.Command{
	Use:   "export <command> [flags]",
	Short: "Export results of a command to a file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// The first argument is the subcommand to execute (e.g., "list", "status")
		subcommandName := args[0]
		subcommandArgs := args[1:]

		// Find the subcommand to execute
		subCmd, _, err := rootCmd.Find(append([]string{subcommandName}, subcommandArgs...))
		if err != nil || subCmd == nil {
			log.Fatalf("Unknown command: %s", subcommandName)
		}

		// Create a buffer to capture the subcommand's output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the subcommand
		subCmd.SetArgs(subcommandArgs)
		subCmd.Execute()

		w.Close()
		output, _ := io.ReadAll(r)
		os.Stdout = oldStdout // Restore original stdout

		// Determine the output format
		var formattedOutput []byte
		switch exportOutputFormat {
		case "json":
			// Assuming the subcommand output is already JSON or can be converted
			// For now, we'll just pass it through. In a real scenario, you'd parse and re-marshal.
			formattedOutput = output
		case "yaml":
			// Assuming the subcommand output is already YAML or can be converted
			// For now, we'll just pass it through.
			formattedOutput = output
		case "table", "":
			formattedOutput = output
		default:
			log.Fatalf("Unsupported output format: %s", exportOutputFormat)
		}

		// Write the captured output to the file
		f, err := os.Create(exportFilePath)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", exportFilePath, err)
		}
		defer f.Close()

		_, err = f.Write(formattedOutput)
		if err != nil {
			log.Fatalf("Error writing to file %s: %v", exportFilePath, err)
		}

		fmt.Printf("Command output exported to %s in %s format.\n", exportFilePath, exportOutputFormat)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	exportCmd.Flags().StringVarP(&exportFilePath, "file", "f", "", "Output file path (required)")
	exportCmd.MarkFlagRequired("file")
}