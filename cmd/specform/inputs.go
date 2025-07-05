package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func LoadInputs(configPath string, inline []string) (map[string]string, error) {
	inputs := map[string]string{}

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read input file: %w", err)
		}
		if err := json.Unmarshal(data, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input JSON: %w", err)
		}
	}

	for _, pair := range inline {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --input format: %s", pair)
		}
		inputs[parts[0]] = parts[1]
	}

	return inputs, nil
}
