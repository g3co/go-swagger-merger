package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMerger(t *testing.T) {
	merger := NewMerger()
	assert.NotNil(t, merger)
	assert.NotNil(t, merger.Swagger)
}

func TestMerger_AddFileYAML(t *testing.T) {
	merger := NewMerger()

	// Test non-existent file
	err := merger.AddFile("nonexistent.yaml")
	assert.Error(t, err)

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
	assert.NoError(t, err)
	defer os.Remove("test.yaml")

	// Test valid file
	err = merger.AddFile("test.yaml")
	assert.NoError(t, err)

	// Verify merged content
	assert.Equal(t, "2.0", merger.Swagger["swagger"])
}

func TestMerger_Save(t *testing.T) {
	merger := NewMerger()
	merger.Swagger["test"] = "value"

	err := merger.Save("test_output.yaml")
	assert.NoError(t, err)
	defer os.Remove("test_output.yaml")

	// Verify file was created
	_, err = os.Stat("test_output.yaml")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	defer os.Remove("test2.yaml")

	// Create second JSON file with overlapping fields
	content3 := []byte(`{
		"swagger": "3.0",
		"info": {
			"title": "Third API",
			"version": "3.0"
		},
		"paths": {
			"/test": {
				"put": {
					"summary": "Test3 PUT, should be created"
				}
			}
		}
	}`)

	err = os.WriteFile("test3.json", content3, 0644)
	assert.NoError(t, err)
	defer os.Remove("test3.json")

	// Add first file
	err = merger.AddFile("test1.yaml")
	assert.NoError(t, err)

	// Verify first file content
	assert.Equal(t, "2.0", merger.Swagger["swagger"])
	assert.Equal(t, "First API", merger.Swagger["info"].(map[string]any)["title"])

	// Add second file
	err = merger.AddFile("test2.yaml")
	assert.NoError(t, err)

	err = merger.AddFile("test3.json")
	assert.NoError(t, err)

	// Verify second file overwrote fields
	assert.Equal(t, "3.0", merger.Swagger["swagger"])
	assert.Equal(t, "Third API", merger.Swagger["info"].(map[string]any)["title"])
	assert.Equal(t, "3.0", merger.Swagger["info"].(map[string]any)["version"])

	// Verify paths were merged
	paths := merger.Swagger["paths"].(map[string]any)
	testPath := paths["/test2"].(map[string]any)
	post := testPath["post"].(map[string]any)
	assert.Equal(t, "Test2 POST, should overwrite first file", post["summary"])

	testPath = paths["/test"].(map[string]any)
	get := testPath["get"].(map[string]any)
	assert.Equal(t, "Test1 GET", get["summary"])

	testPath = paths["/test"].(map[string]any)
	post = testPath["post"].(map[string]any)
	assert.Equal(t, "Test2 POST, should be created, should not overwrite first file", post["summary"])

	testPath = paths["/test"].(map[string]any)
	put := testPath["put"].(map[string]any)
	assert.Equal(t, "Test3 PUT, should be created", put["summary"])
}
