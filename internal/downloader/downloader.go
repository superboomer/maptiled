package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/superboomer/maptiled/internal/loader"
)

// Downloader is a struct for api calls and data save
type Downloader struct {
	url        string
	savePath   string
	httpClient *http.Client

	list *providerList

	setMax bool
}

// NewDownloader create new Download, load all providers
func NewDownloader(url, savePath string, setMax bool) (*Downloader, error) {
	downloader := Downloader{url: url, savePath: savePath, setMax: setMax, httpClient: http.DefaultClient}

	err := os.MkdirAll(savePath, 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}

	list, err := downloader.getProviders()
	if err != nil {
		return nil, fmt.Errorf("error occurred when loading provider list: %w", err)
	}

	downloader.list = list

	return &downloader, nil
}

// GetAllProviders load from list and return as string slice
func (d *Downloader) GetAllProviders() []string {
	var providers = []string{}
	if d.list == nil {
		return providers
	}

	for _, p := range *d.list {
		providers = append(providers, p.Key)
	}

	return providers
}

// DownloadRequest struct which contains all specified data for download tile
type DownloadRequest struct {
	Provider string
	Zoom     int
	Side     int
	Point    *loader.Point
}

// Download download and save specified by DownloadRequest tile
func (d *Downloader) Download(r *DownloadRequest) error {

	req, err := d.createRequest(r)
	if err != nil {
		return fmt.Errorf("error occurred when sending request to the server: err=%w", err)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error occurred when sending request to the server: err=%w", err)
	}

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't readAll body from server answer: err=%w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned invalid status code: code=%d, body=%s", resp.StatusCode, string(img))
	}

	return d.saveImage(r.Point, r.Provider, img)
}

// createRequest build http req
func (d *Downloader) createRequest(r *DownloadRequest) (*http.Request, error) {

	provider, err := d.list.get(r.Provider)
	if err != nil {
		return nil, fmt.Errorf("server dont serve %s provider", r.Provider)
	}

	if provider.MaxZoom < r.Zoom && d.setMax {
		r.Zoom = provider.MaxZoom
	}

	if provider.MaxZoom < r.Zoom {
		return nil, fmt.Errorf("provider %s max zoom %d", r.Provider, provider.MaxZoom)
	}

	return http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/map?provider=%s&lat=%.6f&long=%.6f&zoom=%d&side=%d", d.url, r.Provider, r.Point.Lat, r.Point.Long, r.Zoom, r.Side),
		http.NoBody,
	)
}

// saveImage saves an image file to disk
func (d *Downloader) saveImage(point *loader.Point, provider string, img []byte) error {

	filePath := filepath.Join(d.savePath, fmt.Sprintf("%s_%s_%s.jpeg", point.Name, point.ID, provider))

	file, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write(img)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}

	return nil
}
