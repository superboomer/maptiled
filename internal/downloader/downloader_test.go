package downloader

import (
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
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with dummy image data
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("test image data"))
	}))
	defer server.Close()

	testpath := filepath.Join(os.TempDir(), "test_download")
	err := os.MkdirAll(testpath, 0o700)
	assert.NoError(t, err)

	// Setup downloader with mock client
	downloader := &Downloader{
		httpClient: http.DefaultClient,
		savePath:   testpath, // Ensure this directory exists or adjust accordingly
		url:        server.URL,
		list: &providerList{
			"testProvider": {Name: "Test Provider", Key: "testProvider", MaxZoom: 19},
		},
	}

	// Define download request
	request := &DownloadRequest{
		Provider: "testProvider",
		Zoom:     10,
		Side:     256,
		Point:    &loader.Point{Lat: 12.34, Long: 56.78, Name: "TestPoint", ID: "100"},
	}

	// Call Download method
	err = downloader.Download(request)

	// Assert no error occurred
	assert.NoError(t, err)

	// Verify file was saved correctly
	filePath := filepath.Join(downloader.savePath, fmt.Sprintf("%s_%s_%s.jpeg", request.Point.Name, request.Point.ID, request.Provider))
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "test image data", string(content))

	os.Remove(testpath)
}
