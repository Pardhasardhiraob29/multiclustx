package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
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
		// Find the subcommand to execute
		subCmd, _, err := rootCmd.Find(args)
		if err != nil || subCmd == nil {
			log.Fatalf("Unknown command: %s", args[0])
		}

		// Create a buffer to capture the subcommand's output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the subcommand
		subCmd.SetArgs(args[1:]) // Pass remaining args to subcommand
		subCmd.Execute()

		w.Close()
		output, _ := io.ReadAll(r)
		os.Stdout = oldStdout // Restore original stdout

		// Write the captured output to the file
		f, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", filePath, err)
		}
		defer f.Close()

		// For now, we just write the raw output. 
		// In a real scenario, you'd parse 'output' and format it based on 'outputFormat'.
		_, err = f.Write(output)
		if err != nil {
			log.Fatalf("Error writing to file %s: %v", filePath, err)
		}

		fmt.Printf("Command output exported to %s in raw format.\n", exportFilePath)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportOutputFormat, "output", "o", "json", "Output format (json, yaml, table)")
	exportCmd.Flags().StringVarP(&exportFilePath, "file", "f", "", "Output file path (required)")
	exportCmd.MarkFlagRequired("file")
}
