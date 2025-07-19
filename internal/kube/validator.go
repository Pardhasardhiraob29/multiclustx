package kube

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// PingTest checks the reachability of a Kubernetes cluster's API server.
func PingTest(kubeconfigPath, contextName string) error {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("error loading kubeconfig: %w", err)
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{CurrentContext: contextName})

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("error creating rest config: %w", err)
	}

	// Attempt to create a clientset to check reachability
	_, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to API server: %w", err)
	}

	return nil
}

// GetServerVersion retrieves the Kubernetes API server version.
func GetServerVersion(kubeconfigPath, contextName string) (string, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return "", fmt.Errorf("error loading kubeconfig: %w", err)
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{CurrentContext: contextName})

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return "", fmt.Errorf("error creating rest config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", fmt.Errorf("error creating clientset: %w", err)
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("error getting server version: %w", err)
	}

	return serverVersion.GitVersion, nil
}