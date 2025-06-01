package helpers

import (
	"os"
	"testing"
)

func TestNewMerger(t *testing.T) {
	merger := NewMerger()
	if merger == nil {
		t.Error("Expected non-nil merger")
	}
	if merger.Swagger == nil {
		t.Error("Expected non-nil Swagger map")
	}
}

func TestMerger_AddFile(t *testing.T) {
	merger := NewMerger()

	// Test non-existent file
	err := merger.AddFile("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Create a test YAML file
	content := []byte(`
swagger: "2.0"
info:
  title: Test API
paths:
  /test:
    get:
      summary: Test endpoint
`)

	err = os.WriteFile("test.yaml", content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.yaml")

	// Test valid file
	err = merger.AddFile("test.yaml")
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	// Verify merged content
	if merger.Swagger["swagger"] != "2.0" {
		t.Error("Expected swagger version 2.0")
	}
}

func TestMerger_Save(t *testing.T) {
	merger := NewMerger()
	merger.Swagger["test"] = "value"

	err := merger.Save("test_output.yaml")
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	defer os.Remove("test_output.yaml")

	// Verify file was created
	_, err = os.Stat("test_output.yaml")
	if err != nil {
		t.Error("Expected output file to exist")
	}
}

func TestMerger_MergeMultipleFiles(t *testing.T) {
	merger := NewMerger()

	// Create first YAML file
	content1 := []byte(`
swagger: "2.0"
info:
  title: First API
  version: "1.0"
paths:
  /test:
    get:
      summary: Test1 GET
  /test2:
    post:
      summary: Test1 POST
`)

	err := os.WriteFile("test1.yaml", content1, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test1.yaml")

	// Create second YAML file with overlapping fields
	content2 := []byte(`
swagger: "3.0"
info:
  title: Second API
  version: "2.0"
paths:
  /test:
    post:
      summary: Test2 POST, should be created, should not overwrite first file
  /test2:
    post:
      summary: Test2 POST, should overwrite first file
`)

	err = os.WriteFile("test2.yaml", content2, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test2.yaml")

	// Add first file
	err = merger.AddFile("test1.yaml")
	if err != nil {
		t.Error("Unexpected error adding first file:", err)
	}

	// Verify first file content
	if merger.Swagger["swagger"] != "2.0" {
		t.Error("Expected swagger version 2.0 from first file")
	}
	if merger.Swagger["info"].(map[string]any)["title"] != "First API" {
		t.Error("Expected title from first file")
	}

	// Add second file
	err = merger.AddFile("test2.yaml")
	if err != nil {
		t.Error("Unexpected error adding second file:", err)
	}

	// Verify second file overwrote fields
	if merger.Swagger["swagger"] != "3.0" {
		t.Error("Expected swagger version to be overwritten to 3.0")
	}
	if merger.Swagger["info"].(map[string]interface{})["title"] != "Second API" {
		t.Error("Expected title to be overwritten by second file")
	}
	if merger.Swagger["info"].(map[string]interface{})["version"] != "2.0" {
		t.Error("Expected version to be overwritten to 2.0")
	}

	// Verify paths were merged
	paths := merger.Swagger["paths"].(map[string]interface{})
	testPath := paths["/test2"].(map[string]interface{})
	post := testPath["post"].(map[string]interface{})
	if post["summary"] != "Test2 POST, should overwrite first file" {
		t.Error("Expected test endpoint summary to be overwritten")
	}
}
