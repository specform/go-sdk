package specform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/specform/go-sdk/internal"
)

type CompileOptions struct {
	Strict  bool
	Stdout  bool
	Verbose bool
}

type CompileResult struct {
	Source     string // input file
	OutputPath string // if written to disk
	RawJSON    []byte
}

func CompileSpecFiles(files []string, outputDir string, opts CompileOptions) ([]CompileResult, error) {
	var results []CompileResult

	for _, file := range files {
		compiledPrompt, err := internal.ParseSpecFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", file, err)
		}

		raw, err := json.MarshalIndent(compiledPrompt, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to encode JSON for %s: %w", file, err)
		}

		if opts.Stdout {
			fmt.Println(string(raw))
			results = append(results, CompileResult{
				Source:  file,
				RawJSON: raw,
			})
			continue
		}

		specName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		outPath := filepath.Join(outputDir, specName+".prompt.json")

		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output dir: %w", err)
		}

		if err := os.WriteFile(outPath, raw, 0644); err != nil {
			return nil, fmt.Errorf("failed to write compiled file: %w", err)
		}

		results = append(results, CompileResult{
			Source:     file,
			OutputPath: outPath,
			RawJSON:    raw,
		})
	}

	return results, nil
}
