package downloader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/mtiled/internal/loader"
)

// TestDownload tests the Download function
func TestDownload(t *testing.T) {
	// Test case: Successful download
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("test image data"))
	}))
	defer server.Close()

	testpath := filepath.Join(os.TempDir(), "test_download_success")
	err := os.MkdirAll(testpath, 0o700)
	assert.NoError(t, err)
	defer os.RemoveAll(testpath)

	downloader := &Downloader{
		httpClient: http.DefaultClient,
		savePath:   testpath,
		url:        server.URL,
		list:       &providerList{"testProvider": {Name: "Test Provider", Key: "testProvider", MaxZoom: 19}},
	}

	request := &DownloadRequest{
		Provider: "testProvider",
		Zoom:     10,
		Side:     256,
		Point:    &loader.Point{Lat: 12.34, Long: 56.78, Name: "TestPoint", ID: "100"},
	}

	err = downloader.Download(request)
	assert.NoError(t, err)

	filePath := filepath.Join(testpath, fmt.Sprintf("%s_%s_%s.jpeg", request.Point.Name, request.Point.ID, request.Provider))
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// Test case: Non-existent directory
	testpathNonExistent := filepath.Join(os.TempDir(), "non_existent_directory")
	downloaderNonExistent := &Downloader{
		httpClient: http.DefaultClient,
		savePath:   testpathNonExistent,
		url:        server.URL,
		list:       &providerList{"testProvider": {Name: "Test Provider", Key: "testProvider", MaxZoom: 19}},
	}

	err = downloaderNonExistent.Download(request)
	assert.ErrorContains(t, err, "failed to create file")

	// Test case: Server returns an error status code
	serverError := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverError.Close()

	downloaderServerError := &Downloader{
		httpClient: http.DefaultClient,
		savePath:   testpath,
		url:        serverError.URL,
		list:       &providerList{"testProvider": {Name: "Test Provider", Key: "testProvider", MaxZoom: 19}},
	}

	err = downloaderServerError.Download(request)
	assert.ErrorContains(t, err, "server returned invalid status code")
}

func TestNewDownloader(t *testing.T) {
	// Test case: Successful creation of Downloader
	tempDir := filepath.Join(os.TempDir(), "test_new_downloader_success")
	defer os.RemoveAll(tempDir) // Clean up after test

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

	downloader, err := NewDownloader(server.URL, tempDir, false)
	assert.NoError(t, err)
	assert.NotNil(t, downloader)
	assert.Equal(t, downloader.url, server.URL)
	assert.Equal(t, downloader.savePath, tempDir)
	assert.False(t, downloader.setMax)

	// Test case: Error during directory creation
	_, err = NewDownloader(server.URL, "", false)
	assert.ErrorContains(t, err, "failed to create save directory")
}

func TestGetAllProviders(t *testing.T) {
	// Setup: Create a Downloader instance with a mock provider list
	downloader := &Downloader{
		list: &providerList{
			"testProvider1": {Name: "Test Provider 1", Key: "testProvider1", MaxZoom: 19},
			"testProvider2": {Name: "Test Provider 2", Key: "testProvider2", MaxZoom: 18},
		},
	}

	// Test case: Successfully retrieve all providers
	providers := downloader.GetAllProviders()
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "testProvider1")
	assert.Contains(t, providers, "testProvider2")

	// Test case: Handling an empty provider list
	downloader.list = &providerList{}
	providers = downloader.GetAllProviders()
	assert.Empty(t, providers)
}
