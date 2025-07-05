package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CompileSpecFile(path string, outputDir string) (string, error) {
	scenario, err := ParseSpecFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse spec file: %w", err)
	}

	// Get the filename without the extension
	specName := strings.TrimSuffix((filepath.Base(path)), filepath.Ext(path))
	outFile := filepath.Join(outputDir, specName+".prompt.json")

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the output file in the output directory
	f, err := os.Create(outFile)
	if err != nil {
		return "", fmt.Errorf("failed to create compiled spec file: %w", err)
	}

	// Ensure the file is closed after writing
	defer f.Close()

	// Create a JSON encoder and set indentation
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	// Encode the scenario struct to JSON
	if err := e.Encode(scenario); err != nil {
		return "", fmt.Errorf("failed to encode spec as JSON: %w", err)
	}

	return outFile, nil
}
