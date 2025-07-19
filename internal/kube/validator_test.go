package kube

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestPingTest(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	tmpDir, err := os.MkdirTemp("", "kubeconfig-test")
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
	if err := os.WriteFile(testKubeconfigPath, []byte(kubeconfigContent), 0600); err != nil {
		t.Fatalf("Failed to write kubeconfig file: %v", err)
	}

	// Test failed ping (connection refused is expected for localhost:8080)
	err = PingTest(testKubeconfigPath, "test-context")
	if err == nil {
		t.Error("PingTest did not return an error for unreachable server")
	}

	// Test failed ping (e.g., invalid kubeconfig path)
	err = PingTest("/nonexistent/path/kubeconfig", "test-context")
	if err == nil {
		t.Error("PingTest did not return an error for invalid kubeconfig path")
	}
}

func TestGetServerVersion(t *testing.T) {
	// Create a temporary kubeconfig file for testing
	tmpDir, err := os.MkdirTemp("", "kubeconfig-test")
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
	if err := os.WriteFile(testKubeconfigPath, []byte(kubeconfigContent), 0600); err != nil {
		t.Fatalf("Failed to write kubeconfig file: %v", err)
	}

	// Test failed version retrieval (connection refused is expected for localhost:8080)
	_, err = GetServerVersion(testKubeconfigPath, "test-context")
	if err == nil {
		t.Error("GetServerVersion did not return an error for unreachable server")
	}

	// Test failed version retrieval (e.g., invalid kubeconfig path)
	_, err = GetServerVersion("/nonexistent/path/kubeconfig", "test-context")
	if err == nil {
		t.Error("GetServerVersion did not return an error for invalid kubeconfig path")
	}
}