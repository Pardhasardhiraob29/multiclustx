package kube

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd/api"
)

func TestLoadKubeconfig(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	tmpDir, err := ioutil.TempDir("", "kubeconfig-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testKubeconfigPath := filepath.Join(tmpDir, "config")
	// Minimal valid kubeconfig content
	kubeconfigContent := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    token: some-token
`
	if err := ioutil.WriteFile(testKubeconfigPath, []byte(kubeconfigContent), 0600); err != nil {
		t.Fatalf("Failed to write kubeconfig file: %v", err)
	}

	// Test loading with explicit path
	config, err := LoadKubeconfig(testKubeconfigPath)
	if err != nil {
		t.Fatalf("LoadKubeconfig failed: %v", err)
	}

	if config.CurrentContext != "test-context" {
		t.Errorf("Expected current context 'test-context', got '%s'", config.CurrentContext)
	}

	// Test loading with default path (by setting KUBECONFIG env var)
	os.Setenv("KUBECONFIG", testKubeconfigPath)
	defer os.Unsetenv("KUBECONFIG")

	config, err = LoadKubeconfig("")
	if err != nil {
		t.Fatalf("LoadKubeconfig with default path failed: %v", err)
	}

	if config.CurrentContext != "test-context" {
		t.Errorf("Expected current context 'test-context', got '%s'", config.CurrentContext)
	}
}
