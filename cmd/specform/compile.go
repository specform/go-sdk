package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/specform/go-sdk/internal"
	"github.com/spf13/cobra"
)

func NewCompileCommand() *cobra.Command {
	var outputDir string
	var watchFlag bool
	var verbose bool
	var stdout bool

	cmd := &cobra.Command{
		Use:   "compile [file]",
		Short: "Compile .spec.md files into .prompt.json files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(verbose)

			logger.Debug("Compiling prompt spec files", "outputDir", outputDir, "watch", watchFlag)

			var files []string

			for _, path := range args {
				info, err := os.Stat(path)
				if err != nil {
					logger.Error("Failed to stat path", "path", path, "error", err)
					fmt.Fprintf(os.Stderr, "❌ Failed to stat path %s: %v\n", path, err)
					continue
				}

				if info.IsDir() && stdout {
					logger.Error("Cannot compile directory to stdout", "path", path)
					fmt.Fprintf(os.Stderr, "❌ Cannot compile directory %s to stdout\n", path)
					return fmt.Errorf("cannot compile directory %s to stdout", path)
				}

				if info.IsDir() {
					logger.Debug("Found directory", "path", path)
					err := filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
						if err != nil {
							logger.Error("Error walking directory", "path", p, "error", err)
							return err
						}

						if !fi.IsDir() && strings.HasSuffix(p, ".spec.md") {
							logger.Debug("Found prompt spec file", "file", p)
							files = append(files, p)
						}
						return nil
					})

					if err != nil {
						logger.Error("Error walking directory", "path", path, "error", err)
						fmt.Fprintf(os.Stderr, "❌ Error walking directory %s: %v\n", path, err)
					}
				} else {
					logger.Debug("Found prompt spec file", "file", path)
					files = append(files, path)
				}
			}

			// write files to stdout if --stdout is set
			if stdout {
				logger.Debug("Standard output mode enabled", "files", files)
				logger.Debug("Outputting to stdout", "files", files)
				for _, file := range files {
					parsed, err := internal.ParseSpecFile(file)

					if err != nil {
						logger.Error("Error compiling file", "file", file, "error", err)
						fmt.Fprintf(os.Stderr, "❌ Error compiling %s: %v\n", file, err)
						continue
					}

					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "  ")

					if err := enc.Encode(parsed); err != nil {
						logger.Error("Error encoding JSON", "error", err)
						fmt.Fprintf(os.Stderr, "❌ Error encoding JSON: %v\n", err)
						continue
					}

					return nil
				}
			}

			// Create output directory if it doesn't exist
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					logger.Error("Failed to create output directory", "error", err)
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			// Compile each .spec.md file
			for _, file := range files {
				logger.Debug("Compiling file", "file", file)
				outFile, err := internal.CompileSpecFile(file, outputDir)
				if err != nil {
					logger.Error("Error compiling file", "file", file, "error", err)
					fmt.Fprintf(os.Stderr, "❌ Error compiling %s: %v\n", file, err)
					continue
				}
				logger.Debug("Compiled file", "file", file, "output", outFile)
				fmt.Printf("✅ compiled %s → %s\n", file, outFile)
			}

			if watchFlag {
				logger.Debug("Watching for changes", "files", files, "outputDir", outputDir)

				return WatchFiles(files, outputDir)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "build", "Directory to write compiled .prompt.json files")
	cmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "Watch for changes and recompile")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().BoolVarP(&stdout, "stdout", "s", false, "Output compiled JSON to stdout")

	return cmd
}
