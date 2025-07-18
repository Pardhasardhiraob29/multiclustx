package scanner

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ScanSecretsForTokens scans Kubernetes secrets for base64 encoded tokens.
func ScanSecretsForTokens(kubeconfigPath, contextName string) ([]string, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error loading kubeconfig: %w", err)
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{CurrentContext: contextName})

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating rest config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating clientset: %w", err)
	}

	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing secrets: %w", err)
	}

	var foundTokens []string
	tokenRegex := regexp.MustCompile(`^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+$`)

	for _, secret := range secrets.Items {
		for key, value := range secret.Data {
			decodedValue, err := base64.StdEncoding.DecodeString(string(value))
			if err != nil {
				// Not a base64 string, skip
				continue
			}

			if tokenRegex.MatchString(string(decodedValue)) {
				foundTokens = append(foundTokens, fmt.Sprintf("Secret: %s, Key: %s, Value: %s", secret.Name, key, string(decodedValue)))
			}
		}
	}

	return foundTokens, nil
}
