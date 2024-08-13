package downloader

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderList_Get(t *testing.T) {
	providers := providerList{
		"testProvider": {Name: "Test Provider", Key: "testProvider", MaxZoom: 19},
	}

	// Test retrieving an existing provider
	p, err := providers.get("testProvider")
	assert.NoError(t, err)
	assert.Equal(t, "Test Provider", p.Name)
	assert.Equal(t, "testProvider", p.Key)
	assert.Equal(t, 19, p.MaxZoom)

	// Test attempting to retrieve a non-existent provider
	_, err = providers.get("nonExistent")
	assert.ErrorContains(t, err, "not found")
}

func TestGetProviders(t *testing.T) {
	// Define a slice of providers to return as JSON
	providers := []provider{
		{Name: "Test Provider 1", Key: "testProvider1", MaxZoom: 19},
		{Name: "Test Provider 2", Key: "testProvider2", MaxZoom: 18},
	}

	// Convert the providers slice to JSON
	jsonData, err := json.Marshal(providers)
	if err != nil {
		t.Fatalf("Failed to marshal providers: %v", err)
	}

	// Create a mock HTTP server that returns our JSON data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer server.Close()

	// Inject the mock server URL into the Downloader
	downloader := &Downloader{url: server.URL, httpClient: http.DefaultClient}

	// Call getProviders with our mock client
	providerList, err := downloader.getProviders()
	if err != nil {
		t.Fatalf("Failed to get providers: %v", err)
	}

	// Verify the provider list
	assert.Len(t, *providerList, len(providers))
	for _, p := range providers {
		assert.Equal(t, (*providerList)[p.Key].Name, p.Name)
		assert.Equal(t, (*providerList)[p.Key].MaxZoom, p.MaxZoom)
	}
}
