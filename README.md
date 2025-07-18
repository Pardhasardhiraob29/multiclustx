# MultiClustX

A powerful CLI tool to manage multiple Kubernetes clusters.

## Key Features

- Manages multiple kubeconfig files or contexts
- Runs kubectl-like operations across multiple clusters
- Supports labeling, filtering, and grouping of clusters
- Checks cluster reachability and health
- Audits contexts and generates reports
- Supports context RBAC validation
- Integrates with CI and testing tools

## Installation

```bash
go install github.com/your-username/multiclustx@latest
```

## Usage Examples

```bash
# List all clusters
multiclustx list

# Check reachability and labels
multiclustx status

# Get pods in all clusters
multiclustx exec get pods --namespace kube-system --all-clusters

# RBAC scan for specific cluster
multiclustx audit rbac --context=prod-east

# Scan secrets for tokens
multiclustx scan secrets

# Generate shell autocompletion
multiclustx completion bash

# Serve web UI
multiclustx serve-web

# GitOps sync (placeholder)
multiclustx gitops

# Context switch (placeholder)
multiclustx switch <context-name>

# ArgoCD-like sync dashboard (future extension)
```