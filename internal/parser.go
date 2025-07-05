package internal

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/specform/go-sdk/types"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func GenerateHash(id string) string {
	h := sha256.New()
	h.Write([]byte(id))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ParseSpecFile(path string) (*types.CompiledPrompt, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Grab our meta data from the spec file's frontmatter
	var meta types.CompiledPrompt
	body, err := frontmatter.Parse(bytes.NewReader(content), &meta)

	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Walk markdown and collect blocks
	source := text.NewReader(body)
	doc := goldmark.New().Parser().Parse(source)
	blocks := map[string]string{}

	// Because we have custom languages defined in our spec files, we need to
	// walk the tree to extract our code fences with their language and content
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if fc, ok := n.(*ast.FencedCodeBlock); ok {
				// Get the language of the code block
				lang := string(fc.Language(body))

				// Get the content of the code block
				var sb strings.Builder
				for i := range fc.Lines().Len() {
					line := fc.Lines().At(i)
					sb.Write(line.Value(body))
				}
				blocks[lang] = sb.String()
			}
		}
		// If we are not entering, we just return
		return ast.WalkContinue, nil
	})

	// Assign values from blocks to scenario
	scenario := &meta
	scenario.ID = strings.ReplaceAll(strings.ToLower(scenario.Scenario), " ", "-")
	scenario.Hash = GenerateHash(scenario.ID)
	scenario.CreatedAt = time.Now()
	scenario.UpdatedAt = time.Now()
	scenario.SourcePath = path

	// Parse the prompt
	if val, ok := blocks["prompt"]; ok {
		scenario.Prompt = val
	} else {
		return nil, fmt.Errorf("No prompt found in spec file")
	}

	// Parse the inputs
	if val, ok := blocks["inputs"]; ok {
		vars, defaults, err := ParseInputBlock(val)

		if err != nil {
			return nil, fmt.Errorf("failed to parse inputs block: %w", err)
		}
		scenario.Inputs = vars
		scenario.Values = defaults
	}

	// Parse assertions
	if val, ok := blocks["assertions"]; ok {
		scenario.Assertions, err = ParseAssertionsBlock(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse assertions block: %w", err)
		}
	}

	// Optional Snapshot
	if val, ok := blocks["output"]; ok {
		scenario.Snapshot = val
	}

	return scenario, nil
}
