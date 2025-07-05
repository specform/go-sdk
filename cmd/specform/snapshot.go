package main

import (
	"fmt"
	"os"
	"path/filepath"

	specform "github.com/specform/go-sdk/pkg"
	"github.com/specform/go-sdk/types"
	"github.com/spf13/cobra"
)

func NewSnapshotCommand() *cobra.Command {
	var promptPath string
	var inputsPath string
	var inlineInputs []string
	var outputPath string
	var snapshotDir string
	var similarityPath string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Assert and save snapshot of prompt execution",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(verbose)

			logger.Info("Creating snapshot", "promptPath", promptPath, "inputsPath", inputsPath, "outputPath", outputPath, "snapshotDir", snapshotDir)
			compiled, err := loadCompiledPrompt(promptPath)

			if err != nil {
				logger.Error("Failed to load prompt", "error", err)
				return fmt.Errorf("failed to load prompt: %w", err)
			}

			inputs, err := LoadInputs(inputsPath, inlineInputs)
			if err != nil {
				logger.Error("Failed to load inputs", "error", err)
				return fmt.Errorf("failed to load inputs: %w", err)
			}

			output, err := os.ReadFile(outputPath)
			if err != nil {
				logger.Error("Failed to read output", "error", err)
				return fmt.Errorf("failed to read output: %w", err)
			}

			simScores, err := LoadSimilarityScores(similarityPath)
			if err != nil {
				logger.Error("Failed to load similarity scores", "error", err)
				return fmt.Errorf("failed to load similarity scores: %w", err)
			}

			ctx := &types.AssertionContext{
				SemanticScores: simScores,
			}

			results := specform.RunAssertions(string(output), compiled.Assertions, ctx)
			passed := true
			for _, r := range results {
				logger.Debug("Assertion result", "message", r.Message, "passed", r.Passed)
				fmt.Printf("%s\n", r.Message)

				if !r.Passed {
					passed = false
				}
			}

			if !passed {
				logger.Info("Some assertions failed. Snapshot not saved")
				fmt.Println("❌ Some assertions failed. Snapshot not saved.")
				os.Exit(1)
			}

			if err := os.MkdirAll(snapshotDir, 0755); err != nil {
				logger.Error("Failed to create snapshot directory", "error", err)
				return fmt.Errorf("failed to create snapshot dir: %w", err)
			}

			snapshotPath := filepath.Join(snapshotDir, compiled.ID+".snap.json")
			err = specform.SaveSnapshot(snapshotPath, compiled, string(output), results, inputs)
			if err != nil {
				logger.Error("Failed to save snapshot", "error", err)
				return fmt.Errorf("failed to save snapshot: %w", err)
			}

			logger.Info("Snapshot saved", "path", snapshotPath)
			fmt.Printf("✅ Snapshot saved to %s\n", snapshotPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&promptPath, "prompt", "", "Path to compiled .prompt.json")
	cmd.Flags().StringVar(&inputsPath, "inputs", "", "Path to inputs.json")
	cmd.Flags().StringArrayVar(&inlineInputs, "input", nil, "Inline input as key=value")
	cmd.Flags().StringVar(&outputPath, "output", "", "Path to LLM output.txt")
	cmd.Flags().StringVar(&snapshotDir, "out", "snapshots", "Directory to save snapshots")
	cmd.Flags().StringVar(&similarityPath, "similarity", "", "Optional similarity score JSON file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	_ = cmd.MarkFlagRequired("prompt")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
