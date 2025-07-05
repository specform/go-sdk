package specform

import (
	"testing"

	"github.com/specform/go-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRenderCompilePrompt(t *testing.T) {
	scenario := &types.CompiledPrompt{
		ID:     "test-scenario",
		Prompt: "Summarize this: {{article}} using a {{tone}} tone.",
		Inputs: []string{"article", "tone"},
		Values: map[string]string{
			"article": "Prompt engineering is emerging as a key skill.",
			"tone":    "casual",
		},
	}

	prompt, err := RenderPrompt(scenario, nil, &RenderOptions{
		Strict: true,
	})
	require.NoError(t, err)
	require.Contains(t, prompt, "Prompt engineering is emerging")
	require.Contains(t, prompt, "casual tone")
}

func TestRenderCompiledPrompt_MissingInput(t *testing.T) {
	scenario := &types.CompiledPrompt{
		ID:     "test-scenario",
		Prompt: "Summarize this: {{article}} using a {{tone}} tone.",
		Inputs: []string{"article", "tone"},
		Values: map[string]string{
			"article": "Prompt engineering is emerging as a key skill.",
			// "tone" is missing
		},
	}

	_, err := RenderPrompt(scenario, nil, &RenderOptions{Strict: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required inputs: tone")
}
