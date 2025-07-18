package kube

import (
	"multiclustx/pkg/types"

	"k8s.io/client-go/tools/clientcmd/api"
)

// GetContexts returns a slice of ContextInfo from a kubeconfig.
func GetContexts(config *api.Config) []types.ContextInfo {
	contexts := []types.ContextInfo{}
	for name, context := range config.Contexts {
		contexts = append(contexts, types.ContextInfo{
			Name:     name,
			Cluster:  context.Cluster,
			AuthInfo: context.AuthInfo,
			Namespace: context.Namespace,
		})
	}
	return contexts
}
