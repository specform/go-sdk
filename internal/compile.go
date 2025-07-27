package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func CompileSpecFile(path string, outputDir string) (string, error) {
	compiledPrompt, err := ParseSpecFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse spec file: %w", err)
	}

	// Create the output directory if it doesn't exist
	dirPath := filepath.Join(outputDir, compiledPrompt.Slug)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// create the filename, we're using the file hash here to track versions
	outFileName := fmt.Sprintf("%s.prompt.json", compiledPrompt.Hash[:6])
	outputPath := filepath.Join(dirPath, outFileName)

	// Create the output file in the output directory
	f, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create compiled spec file: %w", err)
	}

	// Ensure the file is closed after writing
	defer f.Close()

	// Create a JSON encoder and set indentation
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	// Encode the scenario struct to JSON
	if err := e.Encode(compiledPrompt); err != nil {
		return "", fmt.Errorf("failed to encode spec as JSON: %w", err)
	}

	return outputPath, nil
}
