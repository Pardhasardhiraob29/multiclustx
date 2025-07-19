package kube

import (
	"fmt"

	"k8s.io/client-go/discovery"
)

// PingTest checks the reachability of a Kubernetes cluster's API server.
func PingTest(discoveryClient discovery.DiscoveryInterface) error {
	// Attempt to create a clientset to check reachability
	_, err := discoveryClient.ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to API server: %w", err)
	}

	return nil
}

// GetServerVersion retrieves the Kubernetes API server version.
func GetServerVersion(discoveryClient discovery.DiscoveryInterface) (string, error) {
	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return "", fmt.Errorf("error getting server version: %w", err)
	}

	return serverVersion.GitVersion, nil
}