package chipper

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewDocumentFromJSON(t *testing.T) {
	// Read the JSON data from the file
	jsonData, err := os.ReadFile("response.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Create a new Document object from the JSON data
	var data []interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	doc := NewDocument(data)

	// Perform assertions on the Document object
	if doc == nil {
		t.Error("Expected non-nil Document object")
	}
	// Add more assertions as needed
}
