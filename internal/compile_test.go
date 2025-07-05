package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompileSpecFile_WritesOutput(t *testing.T) {
	tempDir := t.TempDir()

	outputPath, err := CompileSpecFile("../../../examples/summarize-min.spec.md", tempDir)
	require.NoError(t, err)
	require.FileExists(t, outputPath)

	require.True(t, strings.HasSuffix(outputPath, ".prompt.json"), "expected output filename to end with .prompt.json")

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.Contains(t, string(data), "compiledPrompt")
	require.Contains(t, string(data), "Summarize a technical article")
}
