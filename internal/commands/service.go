package commands

import (
	"fmt"

	"github.com/superboomer/mtiled/internal/downloader"
	"github.com/superboomer/mtiled/internal/loader"
	"github.com/superboomer/mtiled/internal/options"
)

type service struct {
	downloader *downloader.Downloader
	providers  []string
	points     []loader.Point
	zoom       int
	side       int
}

func createService(opts *options.Opts) (*service, error) {

	l := loader.DataLoader{Path: opts.Points}

	points, err := l.Load()
	if err != nil {
		return nil, fmt.Errorf("error occurred when load points: %w", err)
	}

	downloader, err := downloader.NewDownloader(opts.URL, opts.SavePath, opts.SetMax)
	if err != nil {
		return nil, fmt.Errorf("error occurred when init downloader: %w", err)
	}

	var s = &service{
		downloader: downloader,
		points:     points,
		zoom:       opts.Zoom,
		side:       opts.Side,
	}

	if len(opts.Providers) == 0 {
		s.providers = downloader.GetAllProviders()
	} else {
		for _, p := range opts.Providers {
		CHECK:
			for _, a := range downloader.GetAllProviders() {
				if a == p {
					s.providers = append(s.providers, a)
					break CHECK
				}
			}
		}
	}

	if len(s.providers) == 0 {
		return nil, fmt.Errorf("providers not valid")
	}

	return s, nil
}
