package types

// ContextInfo holds information about a single Kubernetes context.
type ContextInfo struct {
	Name      string `json:"name" yaml:"name"`
	Cluster   string `json:"cluster" yaml:"cluster"`
	AuthInfo  string `json:"user" yaml:"user"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}