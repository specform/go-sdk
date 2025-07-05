package main

import (
	"encoding/json"
	"fmt"
	"os"

	specform "github.com/specform/go-sdk/pkg"
	"github.com/specform/go-sdk/types"
	"github.com/spf13/cobra"
)

func NewRenderCommand() *cobra.Command {
	var promptPath string
	var inputsPath string
	var inlineInputs []string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a prompt using a compiled prompt spec and inputs",
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt, err := loadCompiledPrompt(promptPath)
			if err != nil {
				return fmt.Errorf("failed to load prompt: %w", err)
			}

			inputs, err := LoadInputs(inputsPath, inlineInputs)
			if err != nil {
				return fmt.Errorf("failed to load inputs: %w", err)
			}

			rendered, err := specform.RenderPrompt(prompt, inputs, &specform.RenderOptions{Strict: true})
			if err != nil {
				return fmt.Errorf("failed to render: %w", err)
			}

			fmt.Println(rendered)
			return nil
		},
	}

	cmd.Flags().StringVar(&promptPath, "prompt", "", "Path to compiled .prompt.json")
	cmd.Flags().StringVar(&inputsPath, "inputs", "", "Path to inputs.json")
	cmd.Flags().StringArrayVar(&inlineInputs, "input", nil, "Inline input as key=value")

	_ = cmd.MarkFlagRequired("prompt")

	return cmd
}

func loadCompiledPrompt(path string) (*types.CompiledPrompt, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var scenario types.CompiledPrompt
	dec := json.NewDecoder(f)
	if err := dec.Decode(&scenario); err != nil {
		return nil, err
	}
	return &scenario, nil
}
