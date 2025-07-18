package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// LoadKubeconfig loads the kubeconfig from the default path or a specified path.
func LoadKubeconfig(path string) (*api.Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	return config, nil
}
