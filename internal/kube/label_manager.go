package kube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"multiclustx/pkg/types"
)

var labelsFileName = ".multiclustx_labels.json"

// LabelManager handles loading and saving of context labels.
type LabelManager struct {
	labels map[string]map[string]string // map[contextName]map[labelKey]labelValue
	filePath string
}

// NewLabelManager creates a new LabelManager and loads existing labels.
func NewLabelManager() (*LabelManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %w", err)
	}
	filePath := filepath.Join(home, labelsFileName)

	lm := &LabelManager{
		labels: make(map[string]map[string]string),
		filePath: filePath,
	}

	err = lm.loadLabels()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading labels: %w", err)
	}

	return lm, nil
}

// loadLabels loads labels from the file system.
func (lm *LabelManager) loadLabels() error {
	data, err := ioutil.ReadFile(lm.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &lm.labels)
}

// SaveLabels saves labels to the file system.
func (lm *LabelManager) SaveLabels() error {
	data, err := json.MarshalIndent(lm.labels, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling labels: %w", err)
	}

	return ioutil.WriteFile(lm.filePath, data, 0644)
}

// GetLabels returns all labels for a given context.
func (lm *LabelManager) GetLabels(contextName string) map[string]string {
	if l, ok := lm.labels[contextName]; ok {
		return l
	}
	return make(map[string]string)
}

// SetLabel sets a label for a given context.
func (lm *LabelManager) SetLabel(contextName, key, value string) {
	if _, ok := lm.labels[contextName]; !ok {
		lm.labels[contextName] = make(map[string]string)
	}
	lm.labels[contextName][key] = value
}

// DeleteLabel deletes a label from a given context.
func (lm *LabelManager) DeleteLabel(contextName, key string) {
	if l, ok := lm.labels[contextName]; ok {
		delete(l, key)
		if len(l) == 0 {
			delete(lm.labels, contextName)
		}
	}
}

// GetAllContextLabels returns all labels for all contexts.
func (lm *LabelManager) GetAllContextLabels() map[string]map[string]string {
	return lm.labels
}

// FilterContextsByLabel filters a list of contexts based on a given label.
func FilterContextsByLabel(contexts []types.ContextInfo, lm *LabelManager, labelFilter string) []types.ContextInfo {
	filtered := []types.ContextInfo{}
	for _, ctx := range contexts {
		labels := lm.GetLabels(ctx.Name)
		if len(labels) > 0 {
			// Check if the labelFilter matches any of the context's labels
			for key, value := range labels {
				if labelFilter == fmt.Sprintf("%s=%s", key, value) || labelFilter == key {
					filtered = append(filtered, ctx)
					break
				}
			}
		}
	}
	return filtered
}
