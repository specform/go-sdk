package main

import (
	"fmt"
	"os"

	specform "github.com/specform/go-sdk/pkg"
	"github.com/specform/go-sdk/types"
	"github.com/spf13/cobra"
)

func NewTestCommand() *cobra.Command {
	var promptPath string
	var inputsPath string
	var inlineInputs []string
	var outputPath string
	var similarityPath string

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run assertions on a compiled prompt and LLM output",
		RunE: func(cmd *cobra.Command, args []string) error {
			compiled, err := loadCompiledPrompt(promptPath)
			if err != nil {
				return fmt.Errorf("failed to load prompt: %w", err)
			}

			_, err = LoadInputs(inputsPath, inlineInputs)
			if err != nil {
				return fmt.Errorf("failed to load inputs: %w", err)
			}

			output, err := os.ReadFile(outputPath)
			if err != nil {
				return fmt.Errorf("failed to read output: %w", err)
			}

			simScores, err := LoadSimilarityScores(similarityPath)
			if err != nil {
				return fmt.Errorf("failed to load similarity scores: %w", err)
			}

			ctx := &types.AssertionContext{
				SemanticScores: simScores,
			}

			results := specform.RunAssertions(string(output), compiled.Assertions, ctx)
			passed := true
			for _, r := range results {
				fmt.Printf("%s\n", r.Message)
				if !r.Passed {
					passed = false
				}
			}

			if passed {
				fmt.Println("✅ All assertions passed!")
			} else {
				fmt.Println("❌ Some assertions failed")
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&promptPath, "prompt", "", "Path to compiled .prompt.json")
	cmd.Flags().StringVar(&inputsPath, "inputs", "", "Path to inputs.json")
	cmd.Flags().StringArrayVar(&inlineInputs, "input", nil, "Inline input as key=value")
	cmd.Flags().StringVar(&outputPath, "output", "", "Path to LLM output.txt")
	cmd.Flags().StringVar(&similarityPath, "similarity", "", "Optional similarity score JSON file")
	_ = cmd.MarkFlagRequired("prompt")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
