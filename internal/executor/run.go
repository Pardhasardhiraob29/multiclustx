package executor

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecuteKubectlCommand executes a kubectl command against a specific kubeconfig and context.
func ExecuteKubectlCommand(kubeconfigPath, contextName string, args []string) (string, string, error) {
	cmdArgs := []string{"--kubeconfig", kubeconfigPath, "--context", contextName}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("kubectl", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("error executing kubectl command: %w, stderr: %s", err, stderr.String())
	}

	return stdout.String(), stderr.String(), nil
}
