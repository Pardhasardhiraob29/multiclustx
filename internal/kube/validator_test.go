package kube

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestPingTest(t *testing.T) {
	// Mock Kubernetes API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a temporary kubeconfig file pointing to the mock server
	tmpKubeconfig := createTempKubeconfig(t, server.URL)
	defer deleteTempKubeconfig(t, tmpKubeconfig)

	// Test successful ping
	err := PingTest(tmpKubeconfig, "test-context")
	if err != nil {
		t.Errorf("PingTest failed: %v", err)
	}

	// Test failed ping (e.g., invalid kubeconfig path)
	err = PingTest("/nonexistent/path/kubeconfig", "test-context")
	if err == nil {
		t.Error("PingTest did not return an error for invalid kubeconfig path")
	}
}

func TestGetServerVersion(t *testing.T) {
	// Mock Kubernetes API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/version" {
			w.Write([]byte(`{"gitVersion": "v1.23.4"}`))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	// Create a temporary kubeconfig file pointing to the mock server
	tmpKubeconfig := createTempKubeconfig(t, server.URL)
	defer deleteTempKubeconfig(t, tmpKubeconfig)

	// Test successful version retrieval
	version, err := GetServerVersion(tmpKubeconfig, "test-context")
	if err != nil {
		t.Errorf("GetServerVersion failed: %v", err)
	}

	if version != "v1.23.4" {
		t.Errorf("Expected version v1.23.4, got %s", version)
	}

	// Test failed version retrieval (e.g., invalid kubeconfig path)
	version, err = GetServerVersion("/nonexistent/path/kubeconfig", "test-context")
	if err == nil {
		t.Error("GetServerVersion did not return an error for invalid kubeconfig path")
	}
}

// Helper function to create a temporary kubeconfig file
func createTempKubeconfig(t *testing.T, serverURL string) string {
	config := api.Config{
		Clusters: map[string]*api.Cluster{
			"test-cluster": {
				Server: serverURL,
			},
		},
		Contexts: map[string]*api.Context{
			"test-context": {
				Cluster: "test-cluster",
				AuthInfo: "test-user",
			},
		},
		CurrentContext: "test-context",
	}

	// Write the kubeconfig to a temporary file
	tmpFile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp kubeconfig file: %v", err)
	}
	defer tmpFile.Close()

	err = clientcmd.WriteToFile(config, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to write kubeconfig to file: %v", err)
	}

	return tmpFile.Name()
}

// Helper function to delete a temporary kubeconfig file
func deleteTempKubeconfig(t *testing.T, path string) {
	if err := os.Remove(path); err != nil {
		t.Errorf("Failed to delete temp kubeconfig file %s: %v", path, err)
	}
}
