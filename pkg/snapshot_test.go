package specform

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/specform/go-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSaveSnapshot(t *testing.T) {
	tempDir := t.TempDir()
	snapshotPath := filepath.Join(tempDir, "example.snap.json")

	scenario := &types.CompiledPrompt{
		ID:        "test-scenario",
		Hash:      "abc123",
		Prompt:    "What is {{topic}}?",
		Inputs:    []string{"topic"},
		Values:    map[string]string{"topic": "AI"},
		Model:     "gpt-4",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	output := "AI stands for Artificial Intelligence."
	inputs := map[string]string{"topic": "AI"}
	results := []types.AssertionResult{
		{Type: "contains", Value: "Artificial Intelligence", Passed: true, Message: "✔ Output contains"},
		{Type: "matches", Value: "AI", Passed: true, Message: "✔ Output matches"},
	}

	err := SaveSnapshot(snapshotPath, scenario, output, results, inputs)
	require.NoError(t, err)
	require.FileExists(t, snapshotPath)

	// Load the snapshot back using the helper
	loaded, err := LoadSnapshot(snapshotPath)
	require.NoError(t, err)
	require.Equal(t, scenario.ID, loaded.ID)
	require.Equal(t, scenario.Hash, loaded.Hash)
	require.Equal(t, output, loaded.Output)
	require.True(t, loaded.Passed)
	require.Len(t, loaded.Assertions, 2)
}
