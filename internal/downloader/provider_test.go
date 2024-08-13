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
	assert.NoError(t, err)

	// Test case: Successful retrieval of providers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer server.Close()

	downloader := &Downloader{url: server.URL, httpClient: http.DefaultClient}
	providerList, err := downloader.getProviders()
	assert.NoError(t, err)
	assert.Len(t, *providerList, len(providers))
	for _, p := range providers {
		assert.Equal(t, (*providerList)[p.Key].Name, p.Name)
		assert.Equal(t, (*providerList)[p.Key].MaxZoom, p.MaxZoom)
	}

	// Test case: Server returns an empty response
	serverEmpty := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverEmpty.Close()

	downloaderEmpty := &Downloader{url: serverEmpty.URL, httpClient: http.DefaultClient}
	_, err = downloaderEmpty.getProviders()
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	// Test case: Server returns an error status code
	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	downloaderError := &Downloader{url: serverError.URL, httpClient: http.DefaultClient}
	_, err = downloaderError.getProviders()
	assert.ErrorContains(t, err, "server returned invalid status code")

	// Test case: Malformed JSON response from server
	serverMalformedJSON := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{this is not valid json"))
	}))
	defer serverMalformedJSON.Close()

	downloaderMalformedJSON := &Downloader{url: serverMalformedJSON.URL, httpClient: http.DefaultClient}
	_, err = downloaderMalformedJSON.getProviders()
	assert.ErrorContains(t, err, "invalid character")
}
