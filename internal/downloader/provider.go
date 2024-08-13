package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type provider struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	MaxZoom int    `json:"max_zoom"`
}

type providerList map[string]provider

func (l providerList) get(name string) (*provider, error) {
	p, ok := l[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return &p, nil
}

func (d *Downloader) getProviders() (*providerList, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/provider", d.url),
		http.NoBody,
	)

	if err != nil {
		return nil, fmt.Errorf("error occurred when sending request to the server: err=%w", err)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred when sending request to the server: err=%w", err)
	}

	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't readAll body from server answer: err=%w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned invalid status code: code=%d", resp.StatusCode)
	}

	var providers = []provider{}

	err = json.Unmarshal(byt, &providers)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshal result json: %w", err)
	}

	var list = providerList{}

	for _, p := range providers {
		list[p.Key] = p
	}

	return &list, nil
}
