package kube

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

// MockDiscoveryClient is a mock implementation of DiscoveryInterface.
type MockDiscoveryClient struct {
	mock.Mock
}

func (m *MockDiscoveryClient) ServerVersion() (*version.Info, error) {
	args := m.Called()
	return args.Get(0).(*version.Info), args.Error(1)
}

func (m *MockDiscoveryClient) RESTClient() rest.Interface {
	args := m.Called()
	return args.Get(0).(rest.Interface)
}

// MockClientset is a mock implementation of kubernetes.Interface.
type MockClientset struct {
	mock.Mock
}

func (m *MockClientset) Discovery() discovery.DiscoveryInterface {
	args := m.Called()
	return args.Get(0).(discovery.DiscoveryInterface)
}

func (m *MockClientset) CoreV1() corev1typed.CoreV1Interface {
	args := m.Called()
	return args.Get(0).(corev1typed.CoreV1Interface)
}

// ... (other mock methods for other API groups if needed)

func TestPingTest(t *testing.T) {
	// Create a mock clientset
	mockClientset := new(MockClientset)
	mockDiscoveryClient := new(MockDiscoveryClient)

	// Configure mock for successful ping
	mockClientset.On("Discovery").Return(mockDiscoveryClient)
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{}, nil)

	// Create a dummy kubeconfig for the test
	config := &api.Config{
		CurrentContext: "test-context",
		Contexts: map[string]*api.Context{
			"test-context": {},
		},
	}
	// Write the kubeconfig to a temporary file
	tmpFile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp kubeconfig file: %v", err)
	}
	defer tmpFile.Close()

	err = clientcmd.WriteToFile(*config, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to write kubeconfig to file: %v", err)
	}

	// Test successful ping
	err = PingTest(tmpFile.Name(), "test-context")
	if err != nil {
		t.Errorf("PingTest failed: %v", err)
	}

	// Configure mock for failed ping
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{}, errors.New("connection refused"))

	// Test failed ping
	err = PingTest(tmpFile.Name(), "test-context")
	if err == nil {
		t.Error("PingTest did not return an error for failed connection")
	}
}

func TestGetServerVersion(t *testing.T) {
	// Create a mock clientset
	mockClientset := new(MockClientset)
	mockDiscoveryClient := new(MockDiscoveryClient)

	// Configure mock for successful version retrieval
	mockClientset.On("Discovery").Return(mockDiscoveryClient)
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{GitVersion: "v1.23.4"}, nil)

	// Create a dummy kubeconfig for the test
	config := &api.Config{
		CurrentContext: "test-context",
		Contexts: map[string]*api.Context{
			"test-context": {},
		},
	}
	// Write the kubeconfig to a temporary file
	tmpFile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		t.Fatalf("Failed to create temp kubeconfig file: %v", err)
	}
	defer tmpFile.Close()

	err = clientcmd.WriteToFile(*config, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to write kubeconfig to file: %v", err)
	}

	// Test successful version retrieval
	version, err := GetServerVersion(tmpFile.Name(), "test-context")
	if err != nil {
		t.Errorf("GetServerVersion failed: %v", err)
	}

	if version != "v1.23.4" {
		t.Errorf("Expected version v1.23.4, got %s", version)
	}

	// Configure mock for failed version retrieval
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{}, errors.New("failed to get version"))

	// Test failed version retrieval
	_, err = GetServerVersion(tmpFile.Name(), "test-context")
	if err == nil {
		t.Error("GetServerVersion did not return an error for failed version retrieval")
	}
}