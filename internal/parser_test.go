package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserSpecFile(t *testing.T) {
	spec, err := ParseSpecFile("../../../examples/summarize-min.spec.md")
	require.NoError(t, err)
	require.NotNil(t, spec)

	// Meta Parsing
	require.Equal(t, "Summarize a technical article", spec.Scenario)
	require.Equal(t, "summarize-a-technical-article", spec.ID)
	require.Equal(t, GenerateHash("summarize-a-technical-article"), spec.Hash)

	require.Equal(t, "gpt-4", spec.Model)
	require.Equal(t, 0.3, spec.Temperature)
	require.Equal(t, []string{"summarization", "test"}, spec.Tags)

	// Prompt Parsing
	require.Contains(t, spec.Prompt, "{{article}}")
	require.Contains(t, spec.Prompt, "{{tone}}")

	// Input Parsing
	require.ElementsMatch(t, []string{"article", "tone"}, spec.Inputs)
	require.Contains(t, spec.Inputs, "article")
	require.Equal(t, "casual", spec.Values["tone"])
	require.Contains(t, spec.Values["article"], "Webhooks enable real-time communication")

	// Assertion Parsing
	require.Len(t, spec.Assertions, 3)
	require.Equal(t, "contains", spec.Assertions[0].Type)
	require.Equal(t, "real time", spec.Assertions[0].Value)
	require.Equal(t, "matches", spec.Assertions[1].Type)
	require.Equal(t, "/HTTP/i", spec.Assertions[1].Value)
	require.Equal(t, "semantic-similarity", spec.Assertions[2].Type)
	require.Equal(t, "event-driven communication", spec.Assertions[2].Value)

	require.NotEmpty(t, spec.CreatedAt)
	require.NotEmpty(t, spec.UpdatedAt)
}

func TestParseSpec_MissingPromptFails(t *testing.T) {
	_, err := ParseSpecFile("../examples/prompts/missing_prompt.spec.md")
	require.Error(t, err)
	require.Contains(t, err.Error(), "No prompt found in spec file")
}

func TestParseSpec_InvalidInputs(t *testing.T) {
	spec, err := ParseSpecFile("../examples/prompts/invalid_inputs.spec.md")
	require.NoError(t, err)
	require.NotContains(t, spec.Inputs, "article") // failed to parse as key=value
}

func TestParseSpec_InvalidFrontmatter(t *testing.T) {
	_, err := ParseSpecFile("../examples/prompts/invalid_frontmatter.spec.md")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse frontmatter")
}
