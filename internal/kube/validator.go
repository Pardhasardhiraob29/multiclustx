package kube

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// PingTest checks the reachability of a Kubernetes cluster's API server.
func PingTest(config *rest.Config) error {
	// Attempt to create a clientset to check reachability
	_, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to connect to API server: %w", err)
	}

	return nil
}

// GetServerVersion retrieves the Kubernetes API server version.
func GetServerVersion(config *rest.Config) (string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("error creating clientset: %w", err)
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("error getting server version: %w", err)
	}

	return serverVersion.GitVersion, nil
}