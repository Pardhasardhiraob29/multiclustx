package rbac

import (
	"context"
	"fmt"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// CheckRBAC performs a SelfSubjectRulesReview for the given context.
func CheckRBAC(kubeconfigPath, contextName string) (*authorizationv1.SelfSubjectRulesReview, error) {
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

	selfSubjectRulesReview := &authorizationv1.SelfSubjectRulesReview{
		Spec: authorizationv1.SelfSubjectRulesReviewSpec{
			Namespace: "default", // You might want to make this configurable or iterate through namespaces
		},
	}

	response, err := clientset.AuthorizationV1().SelfSubjectRulesReviews().Create(context.TODO(), selfSubjectRulesReview, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("error performing SelfSubjectRulesReview: %w", err)
	}

	return response, nil
}
