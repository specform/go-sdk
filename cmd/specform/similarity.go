package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadSimilarityScores(path string) (map[string]float64, error) {
	if path == "" {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open similarity file: %w", err)
	}
	defer f.Close()

	scores := make(map[string]float64)
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&scores); err != nil {
		return nil, fmt.Errorf("failed to decode similarity file: %w", err)
	}

	return scores, nil
}
