package internal

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/specform/go-sdk/types"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Helper function generates a SHA256 hash from the passed byte slice.
// This is used to create a unique identifier for the compiled prompt.
func GenerateHash(data []byte) string {
	hashed := sha256.Sum256(data)
	return fmt.Sprintf("%x", hashed)
}

// Helper function generates a id based on the slug + file hash
func IDFrom(slug, hash string) string {
	return fmt.Sprintf("%s-%s", slug, hash[:6])
}

// Helper function generates a slug based on the spec's title
func Slugify(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "") // remove non-word characters
	slug = strings.ReplaceAll(slug, " ", "-")                        // replace spaces with dashes
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")      // collapse multiple dashes

	return slug
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

	// Assign values from blocks to compiledPrompt
	compiledPrompt := &meta
	compiledPrompt.Slug = Slugify(compiledPrompt.Title)
	compiledPrompt.Hash = GenerateHash(content)
	compiledPrompt.ID = IDFrom(compiledPrompt.Slug, compiledPrompt.Hash)
	compiledPrompt.CreatedAt = time.Now()
	compiledPrompt.UpdatedAt = time.Now()
	compiledPrompt.SourcePath = path

	// Parse the prompt
	if val, ok := blocks["prompt"]; ok {
		compiledPrompt.Prompt = val
	} else {
		return nil, fmt.Errorf("no prompt found in spec file")
	}

	// Parse the inputs
	if val, ok := blocks["inputs"]; ok {
		vars, defaults, err := ParseInputBlock(val)

		if err != nil {
			return nil, fmt.Errorf("failed to parse inputs block: %w", err)
		}
		compiledPrompt.Inputs = vars
		compiledPrompt.Values = defaults
	}

	// Parse assertions
	if val, ok := blocks["assertions"]; ok {
		compiledPrompt.Assertions, err = ParseAssertionsBlock(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse assertions block: %w", err)
		}
	}

	// Optional Snapshot
	if val, ok := blocks["output"]; ok {
		compiledPrompt.Snapshot = val
	}

	return compiledPrompt, nil
}
