package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// LoadKubeconfig loads the kubeconfig from the given path. If the path is empty,
// it attempts to load from the KUBECONFIG environment variable or the default home directory path.
func LoadKubeconfig(path string) (*api.Config, error) {
	if path != "" {
		config, err := clientcmd.LoadFromFile(path)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// If path is empty, try to load from KUBECONFIG environment variable
	envKubeconfig := os.Getenv("KUBECONFIG")
	if envKubeconfig != "" {
		config, err := clientcmd.LoadFromFile(envKubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// If KUBECONFIG is not set, try the default home directory path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	defaultKubeconfigPath := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.LoadFromFile(defaultKubeconfigPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}