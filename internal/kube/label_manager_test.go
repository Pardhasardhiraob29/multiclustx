package kube

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"multiclustx/pkg/types"
)

func TestLabelManager(t *testing.T) {
	// Create a temporary directory for the labels file
	tmpDir, err := ioutil.TempDir("", "label-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testLabelsFilePath := filepath.Join(tmpDir, ".multiclustx_labels.json")

	// Test NewLabelManager and initial state
	lm, err := NewLabelManager(testLabelsFilePath)
	if err != nil {
		t.Fatalf("NewLabelManager failed: %v", err)
	}

	if len(lm.GetAllContextLabels()) != 0 {
		t.Errorf("Expected no labels initially, got %v", lm.GetAllContextLabels())
	}

	// Test SetLabel
	lm.SetLabel("context1", "env", "dev")
	lm.SetLabel("context1", "region", "us-east-1")
	lm.SetLabel("context2", "env", "prod")

	// Test GetLabels
	labels1 := lm.GetLabels("context1")
	if labels1["env"] != "dev" || labels1["region"] != "us-east-1" || len(labels1) != 2 {
		t.Errorf("Expected labels for context1 to be env=dev, region=us-east-1, got %v", labels1)
	}

	labels2 := lm.GetLabels("context2")
	if labels2["env"] != "prod" || len(labels2) != 1 {
		t.Errorf("Expected labels for context2 to be env=prod, got %v", labels2)
	}

	// Test SaveLabels
	err = lm.SaveLabels()
	if err != nil {
		t.Fatalf("SaveLabels failed: %v", err)
	}

	// Test loading labels from file
	lm2, err := NewLabelManager(testLabelsFilePath)
	if err != nil {
		t.Fatalf("NewLabelManager (after save) failed: %v", err)
	}

	labels1_loaded := lm2.GetLabels("context1")
	if labels1_loaded["env"] != "dev" || labels1_loaded["region"] != "us-east-1" || len(labels1_loaded) != 2 {
		t.Errorf("Expected loaded labels for context1 to be env=dev, region=us-east-1, got %v", labels1_loaded)
	}

	// Test DeleteLabel
	lm2.DeleteLabel("context1", "env")
	labels1_deleted := lm2.GetLabels("context1")
	if labels1_deleted["env"] != "" || labels1_deleted["region"] != "us-east-1" || len(labels1_deleted) != 1 {
		t.Errorf("Expected labels for context1 after delete to be region=us-east-1, got %v", labels1_deleted)
	}

	lm2.DeleteLabel("context1", "region")
	labels1_all_deleted := lm2.GetLabels("context1")
	if len(labels1_all_deleted) != 0 {
		t.Errorf("Expected no labels for context1 after all deletes, got %v", labels1_all_deleted)
	}

	// Test FilterContextsByLabel
	contexts := []types.ContextInfo{
		{Name: "context1"},
		{Name: "context2"},
		{Name: "context3"},
	}

	lm3, err := NewLabelManager(testLabelsFilePath)
	if err != nil {
		t.Fatalf("NewLabelManager (for filter) failed: %v", err)
	}
	lm3.SetLabel("context1", "env", "dev")
	lm3.SetLabel("context2", "env", "prod")
	lm3.SaveLabels()

	filtered := FilterContextsByLabel(contexts, lm3, "env=dev")
	if len(filtered) != 1 || filtered[0].Name != "context1" {
		t.Errorf("Expected 1 context (context1) for filter env=dev, got %v", filtered)
	}

	filtered = FilterContextsByLabel(contexts, lm3, "env")
	if len(filtered) != 2 || (filtered[0].Name != "context1" && filtered[1].Name != "context1") || (filtered[0].Name != "context2" && filtered[1].Name != "context2") {
		t.Errorf("Expected 2 contexts (context1, context2) for filter env, got %v", filtered)
	}
}