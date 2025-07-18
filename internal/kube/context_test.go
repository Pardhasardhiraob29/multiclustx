package kube

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd/api"
)

func TestGetContexts(t *testing.T) {
	// Create a mock api.Config object
	config := &api.Config{
		Contexts: map[string]*api.Context{
			"test-context-1": {
				Cluster:  "test-cluster-1",
				AuthInfo: "test-user-1",
				Namespace: "default",
			},
			"test-context-2": {
				Cluster:  "test-cluster-2",
				AuthInfo: "test-user-2",
				Namespace: "kube-system",
			},
		},
	}

	// Call GetContexts
	contexts := GetContexts(config)

	// Assertions
	if len(contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(contexts))
	}

	// Check context-1
	found1 := false
	for _, ctx := range contexts {
		if ctx.Name == "test-context-1" {
			found1 = true
			if ctx.Cluster != "test-cluster-1" {
				t.Errorf("Expected cluster 'test-cluster-1', got '%s' for context 'test-context-1'", ctx.Cluster)
			}
			if ctx.AuthInfo != "test-user-1" {
				t.Errorf("Expected user 'test-user-1', got '%s' for context 'test-context-1'", ctx.AuthInfo)
			}
			if ctx.Namespace != "default" {
				t.Errorf("Expected namespace 'default', got '%s' for context 'test-context-1'", ctx.Namespace)
			}
			break
		}
	}
	if !found1 {
		t.Errorf("Context 'test-context-1' not found")
	}

	// Check context-2
	found2 := false
	for _, ctx := range contexts {
		if ctx.Name == "test-context-2" {
			found2 = true
			if ctx.Cluster != "test-cluster-2" {
				t.Errorf("Expected cluster 'test-cluster-2', got '%s' for context 'test-context-2'", ctx.Cluster)
			}
			if ctx.AuthInfo != "test-user-2" {
				t.Errorf("Expected user 'test-user-2', got '%s' for context 'test-context-2'", ctx.AuthInfo)
			}
			if ctx.Namespace != "kube-system" {
				t.Errorf("Expected namespace 'kube-system', got '%s' for context 'test-context-2'", ctx.Namespace)
			}
			break
		}
	}
	if !found2 {
		t.Errorf("Context 'test-context-2' not found")
	}
}
