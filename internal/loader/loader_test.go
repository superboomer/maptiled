package loader

import (
	"os"
	"testing"
)

func TestLoad_SuccessfulLoad(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Write valid JSON to the temporary file
	jsonContent := []byte(`[
		{"lat": 40.712776, "long": -74.005974, "name": "New York", "id": "NYC"},
		{"lat": 34.052235, "long": -118.243683, "name": "Los Angeles", "id": "LA"}
	]`)
	if _, wErr := tmpfile.Write(jsonContent); wErr != nil {
		t.Fatal(wErr)
	}
	if cErr := tmpfile.Close(); cErr != nil {
		t.Fatal(cErr)
	}

	// Use DataLoader to load the data
	dataLoader := DataLoader{Path: tmpfile.Name()}
	points, err := dataLoader.Load()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check the loaded data
	if len(points) != 2 {
		t.Errorf("expected 2 points, got %d", len(points))
	}
	if points[0].Name != "New York" || points[0].ID != "NYC" || points[0].Lat != 40.712776 || points[0].Long != -74.005974 {
		t.Errorf("point data does not match expected values")
	}
}

func TestLoad_FileDoesNotExist(t *testing.T) {
	dataLoader := DataLoader{Path: "nonexistent_file.json"}

	_, err := dataLoader.Load()
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tmpfile, err := os.CreateTemp("", "invalid-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Write invalid JSON to the temporary file
	if _, wErr := tmpfile.WriteString("{invalid json}"); wErr != nil {
		t.Fatal(wErr)
	}
	if cErr := tmpfile.Close(); cErr != nil {
		t.Fatal(cErr)
	}

	// Use DataLoader to load the data
	dataLoader := DataLoader{Path: tmpfile.Name()}
	_, err = dataLoader.Load()
	if err == nil {
		t.Error("expected an error, got nil")
	}
}
