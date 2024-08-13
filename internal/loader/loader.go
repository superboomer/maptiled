package loader

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataLoader struct {
	Path string
}

type Point struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
	Name string  `json:"name"`
	ID   string  `json:"id"`
}

func (l *DataLoader) Load() ([]Point, error) {
	byt, err := os.ReadFile(l.Path)
	if err != nil {
		return nil, fmt.Errorf("error occurred while loading data file: %w", err)
	}

	result := []Point{}

	err = json.Unmarshal(byt, &result)
	if err != nil {
		return nil, fmt.Errorf("error occurred while unmarshal data file: %w", err)
	}

	return result, nil
}
