package kube

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	version "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"

	corev1typed "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (m *MockDiscoveryClient) Discovery() discovery.DiscoveryInterface {
	args := m.Called()
	return args.Get(0).(discovery.DiscoveryInterface)
}

func (m *MockDiscoveryClient) RESTClientFor(config *rest.Config) (*rest.RESTClient, error) {
	args := m.Called(config)
	return args.Get(0).(*rest.RESTClient), args.Error(1)
}

func (m *MockDiscoveryClient) ServerGroups() (*metav1.APIGroupList, error) {
	args := m.Called()
	return args.Get(0).(*metav1.APIGroupList), args.Error(1)
}

func (m *MockDiscoveryClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	args := m.Called(groupVersion)
	return args.Get(0).(*metav1.APIResourceList), args.Error(1)
}

func (m *MockDiscoveryClient) ServerResources() ([]*metav1.APIResourceList, error) {
	args := m.Called()
	return args.Get(0).([]*metav1.APIResourceList), args.Error(1)
}

func (m *MockDiscoveryClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	args := m.Called()
	return args.Get(0).([]*metav1.APIResourceList), args.Error(1)
}

func (m *MockDiscoveryClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	args := m.Called()
	return args.Get(0).([]*metav1.APIResourceList), args.Error(1)
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
	// Create a mock discovery client
	mockDiscoveryClient := new(MockDiscoveryClient)

	// Configure mock for successful ping
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{}, nil)

	// Test successful ping
	err := PingTest(mockDiscoveryClient)
	if err != nil {
		t.Errorf("PingTest failed: %v", err)
	}

	// Configure mock for failed ping
	mockDiscoveryClient.On("ServerVersion").Return(nil, errors.New("connection refused"))

	// Test failed ping
	err = PingTest(mockDiscoveryClient)
	if err == nil {
		t.Error("PingTest did not return an error for failed connection")
	}
}

func TestGetServerVersion(t *testing.T) {
	// Create a mock discovery client
	mockDiscoveryClient := new(MockDiscoveryClient)

	// Configure mock for successful version retrieval
	mockDiscoveryClient.On("ServerVersion").Return(&version.Info{GitVersion: "v1.23.4"}, nil)

	// Test successful version retrieval
	version, err := GetServerVersion(mockDiscoveryClient)
	if err != nil {
		t.Errorf("GetServerVersion failed: %v", err)
	}

	if version != "v1.23.4" {
		t.Errorf("Expected version v1.23.4, got %s", version)
	}

	// Configure mock for failed version retrieval
	mockDiscoveryClient.On("ServerVersion").Return(nil, errors.New("failed to get version"))

	// Test failed version retrieval
	_, err = GetServerVersion(mockDiscoveryClient)
	if err == nil {
		t.Error("GetServerVersion did not return an error for failed version retrieval")
	}
}
