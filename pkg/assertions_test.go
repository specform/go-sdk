package specform

import (
	"strings"
	"testing"

	"github.com/specform/go-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRunAssertions(t *testing.T) {
	output := "Webhooks allow real-time HTTP communication between systems."

	err := RegisterAssertion("starts-with", func(val, out string, _ *types.AssertionContext) types.AssertionResult {
		passed := strings.HasPrefix(out, val)
		msg := passFailMsg(passed, "Output starts with '%s'", "Output does not start with '%s'", val)
		return types.AssertionResult{
			Type:    "starts-with",
			Value:   val,
			Passed:  passed,
			Message: msg,
		}
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		assertions []types.Assertion
		extras     *types.AssertionContext
		expected   []bool
	}{
		{
			name:       "contains (pass)",
			assertions: []types.Assertion{{Type: "contains", Value: "real time"}},
			expected:   []bool{true},
		},
		{
			name:       "contains (fail)",
			assertions: []types.Assertion{{Type: "contains", Value: "websocket"}},
			expected:   []bool{false},
		},
		{
			name:       "equals (pass)",
			assertions: []types.Assertion{{Type: "equals", Value: "Webhooks allow real-time HTTP communication between systems."}},
			expected:   []bool{true},
		},
		{
			name:       "equals (fail)",
			assertions: []types.Assertion{{Type: "equals", Value: "something else"}},
			expected:   []bool{false},
		},
		{
			name:       "matches (pass)",
			assertions: []types.Assertion{{Type: "matches", Value: "/HTTP/i"}},
			expected:   []bool{true},
		},
		{
			name:       "matches (fail)",
			assertions: []types.Assertion{{Type: "matches", Value: "/WebSocket/"}},
			expected:   []bool{false},
		},
		{
			name:       "matches (invalid regex)",
			assertions: []types.Assertion{{Type: "matches", Value: "[[invalid"}},
			expected:   []bool{false},
		},
		{
			name:       "semantic-similarity (pass)",
			assertions: []types.Assertion{{Type: "semantic-similarity", Value: "event-driven communication"}},
			extras:     &types.AssertionContext{SemanticScores: map[string]float64{"event-driven communication": 0.92}},
			expected:   []bool{true},
		},
		{
			name:       "semantic-similarity (fail)",
			assertions: []types.Assertion{{Type: "semantic-similarity", Value: "event-driven communication"}},
			extras:     &types.AssertionContext{SemanticScores: map[string]float64{"event-driven communication": 0.6}},
			expected:   []bool{false},
		},
		{
			name:       "unknown assertion type",
			assertions: []types.Assertion{{Type: "unknown-type", Value: "whatever"}},
			expected:   []bool{false},
		},
		{
			name:       "custom assertion (starts-with)",
			assertions: []types.Assertion{{Type: "starts-with", Value: "Webhooks"}},
			expected:   []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := RunAssertions(output, tt.assertions, tt.extras)
			require.Equal(t, len(tt.expected), len(results))
			for i, res := range results {
				require.Equal(t, tt.expected[i], res.Passed, res.Message)
			}
		})
	}
}
