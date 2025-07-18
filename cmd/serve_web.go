package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var serveWebCmd = &cobra.Command{
	Use:   "serve-web",
	Short: "Serve the MultiClustX web UI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting MultiClustX web UI...")
		command := exec.Command("go", "run", "main_web.go")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			log.Fatalf("Error starting web server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveWebCmd)
}
